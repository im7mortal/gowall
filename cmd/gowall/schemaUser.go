package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type vendorOauth struct {
	ID string `bson:"id"`
}

type User struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"-" bson:"password"`
	Email    string        `json:"email" bson:"email"`
	Roles    struct {
		Admin   bson.ObjectId `json:"admin" bson:"admin,omitempty"`
		Account bson.ObjectId `json:"account" bson:"account,omitempty"`
	} `json:"roles" bson:"roles"`

	IsActive             string    `json:"isActive" bson:"isActive,omitempty"`
	TimeCreated          time.Time `json:"timeCreated" bson:"timeCreated"`
	ResetPasswordToken   string    `json:"-" bson:"resetPasswordToken,omitempty"`
	ResetPasswordExpires time.Time `json:"resetPasswordExpires" bson:"resetPasswordExpires,omitempty"`

	Twitter  vendorOauth `json:"twitter" bson:"twitter"`
	Github   vendorOauth `json:"github" bson:"github"`
	Facebook vendorOauth `json:"facebook" bson:"facebook"`
	Google   vendorOauth `json:"google" bson:"google"`
	Tumblr   vendorOauth `json:"tumblr" bson:"tumblr"`
	Search   []string    `json:"search" bson:"search"`
}

func (user *User) canPlayRoleOf(role string) bool {
	if role == "admin" && len(user.Roles.Account.String()) > 0 {
		return true
	}
	if role == "account" && len(user.Roles.Account.String()) > 0 {
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
		panic(err)
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
