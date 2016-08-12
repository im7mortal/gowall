package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
	"encoding/json"
	"sync"
)

type vendorOauth struct {
	ID string `bson:"id"`
}

type User struct {
	ID                   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Username             string        `json:"username" bson:"username"`
	Password             string        `json:"-" bson:"password"`
	Email                string        `json:"email" bson:"email"`
	Roles                struct {
							 Admin   bson.ObjectId `json:"admin" bson:"admin,omitempty"`
							 Account bson.ObjectId `json:"account" bson:"account,omitempty"`
						 } `json:"roles" bson:"roles"`

	IsActive             string    `json:"isActive" bson:"isActive,omitempty"`
	TimeCreated          time.Time `json:"timeCreated" bson:"timeCreated"`
	ResetPasswordToken   string    `json:"-" bson:"resetPasswordToken,omitempty"`
	ResetPasswordExpires time.Time `json:"resetPasswordExpires" bson:"resetPasswordExpires,omitempty"`

	Twitter              vendorOauth `json:"twitter" bson:"twitter"`
	Github               vendorOauth `json:"github" bson:"github"`
	Facebook             vendorOauth `json:"facebook" bson:"facebook"`
	Google               vendorOauth `json:"google" bson:"google"`
	Tumblr               vendorOauth `json:"tumblr" bson:"tumblr"`
	Search               []string    `json:"search" bson:"search"`
}

func (user *User) canPlayRoleOf(role string) bool {
	if role == "admin" && len(user.Roles.Admin.Hex()) > 0 {
		return true
	}
	if role == "account" && len(user.Roles.Account.Hex()) > 0 {
		return true
	}
	return false
}

func (user *User) defaultReturnUrl() (returnUrl string) {
	returnUrl = "/"
	if user.canPlayRoleOf("account") {
		returnUrl = "/account/"
	}
	if user.canPlayRoleOf("admin") {
		returnUrl = "/admin/"
	}
	return
}

func (user *User) setPassword(password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		EXCEPTION(err)
	}
	user.Password = string(hashedPassword)
}

func validateUsername(username *string, r *Response) {
	*username = strings.ToLower(*username)
	if len(*username) == 0 {
		r.ErrFor["username"] = "required"
	} else {
		if ok := rUsername.MatchString(*username); !ok {
			r.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
}

func validateEmail(email *string, r *Response) {
	*email = strings.ToLower(*email)
	if len(*email) == 0 {
		r.ErrFor["email"] = "required"
	} else {
		if ok := rEmail.MatchString(*email); !ok {
			r.ErrFor["email"] = `invalid email format`
		}
	}
}

func validatePassword(password *string, r *Response) {
	if len(*password) == 0 {
		r.ErrFor["password"] = "required"
	} else {
		if len(*password) < 8 {
			r.ErrFor["password"] = `too weak password, at least 8 necessary`
		}
	}
}

var UserUniqueIndex mgo.Index = mgo.Index{
	Key:    []string{"username", "email"},
	Unique: true,
}

var UserIndex mgo.Index = mgo.Index{
	Key: []string{"timeCreated", "twitter.id", "github.id", "facebook.id", "google.id", "search"},
}

func (user *User) login(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Set("public", user.ID.Hex())
	sess.Save()
}

func (user *User) isPasswordOk(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (user *User) changePassword(r *Response) (err error) {

	var body struct {
		Confirm  string `json:"confirm"`
		Password string `json:"newPassword"`
	}
	err = json.NewDecoder(r.c.Request.Body).Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}

	// validate
	if len(body.Password) == 0 {
		r.ErrFor["newPassword"] = "required"
	}
	if len(body.Confirm) == 0 {
		r.ErrFor["confirm"] = "required"
	} else if body.Password != body.Confirm {
		r.Errors = append(r.Errors, "Passwords do not match.")
	}

	if r.HasErrors() {
		err = Err
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()

	user.setPassword(body.Password)

	collection := db.C(USERS)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		r.Errors = append(r.Errors, err.Error())
		err = Err
		return
	}

	return
}

func (user *User) changeIdentity(r *Response) (err error) {

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		IsActive string `json:"isActive"`
	}
	err = json.NewDecoder(r.c.Request.Body).Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}

	validateUsername(&body.Username, r)
	validateEmail(&body.Email, r)

	if r.HasErrors() {
		err = Err
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)

	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{
				"username": body.Username,
			},
			bson.M{
				"email": body.Email,
			},
		},
		"_id": bson.M{"$ne": user.ID},
	}).One(nil)
	if err == nil {
		r.Errors = append(r.Errors, "That username or email already exist.")
		err = Err
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}
	user.Username = body.Username
	user.Email = body.Email
	if len(body.IsActive) != 0 {
		user.IsActive = body.IsActive
	}
	err = collection.UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"username": user.Username,
			"email":    user.Email,
			"isActive": user.IsActive,
			"search":   []string{user.Username, user.Email},
		},
	})
	if err != nil {
		r.Errors = append(r.Errors, err.Error())
		err = Err
		return
	}

	// patchAdmin
	// patchAccount
	updateRoles(db, user)

	return
}

func updateRoles(db *mgo.Database, user *User) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := db.C(ADMINS).Update(
			bson.M{
				"roles.admin": user.ID,
			},
			bson.M{
				"$set": bson.M{
					"user": bson.M{
						"id": user.ID,
						"name": user.Username,
					},
				},
			})
		if err != nil {
			if err != mgo.ErrNotFound {
				EXCEPTION(err)
			}
		}
		wg.Done()
	}()
	go func() {
		err := db.C(ACCOUNTS).Update(
			bson.M{
				"roles.account": user.ID,
			},
			bson.M{
				"$set": bson.M{
					"user": bson.M{
						"id": user.ID,
						"name": user.Username,
					},
				},
			})
		if err != nil {
			if err != mgo.ErrNotFound {
				EXCEPTION(err)
			}
		}
		wg.Done()
	}()
	wg.Wait()

	return
}
