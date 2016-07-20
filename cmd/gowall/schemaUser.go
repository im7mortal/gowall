package main

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
)

type vendorOauth struct {
	ID string `bson:"id"`
}

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `json:"-" bson:"password"`
	Email    string        `bson:"email"`
	Roles    struct {
		Admin   bson.ObjectId `bson:"admin,omitempty"`
		Account bson.ObjectId `bson:"account,omitempty"`
	} `bson:"roles"`

	IsActive             string    `bson:"isActive,omitempty"`
	TimeCreated          time.Time `bson:"timeCreated"`
	ResetPasswordToken   string    `json:"-" bson:"resetPasswordToken,omitempty"`
	ResetPasswordExpires time.Time `bson:"resetPasswordExpires,omitempty"`

	Twitter  vendorOauth `bson:"twitter"`
	Github   vendorOauth `bson:"github"`
	Facebook vendorOauth `bson:"facebook"`
	Google   vendorOauth `bson:"google"`
	Tumblr   vendorOauth `bson:"tumblr"`
	Search   []string    `bson:"search"`
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