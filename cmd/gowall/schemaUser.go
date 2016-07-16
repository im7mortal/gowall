package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"regexp"
)

type vendorOauth struct {
	ID string `bson:"id"`
}

type User struct {
	ID                   bson.ObjectId `bson:"_id"`
	Username             string `bson:"username"`
	Password             string `bson:"password"`
	Email                string `bson:"email"`
	Roles                struct {
							 Admin   mgo.DBRef `bson:"admin,omitempty"`
							 Account mgo.DBRef `bson:"account,omitempty"`
						 } `bson:"roles"`

	IsActive             string `bson:"isActive,omitempty"`
	TimeCreated          time.Time `bson:"timeCreated"`
	ResetPasswordToken   string `bson:"resetPasswordToken,omitempty"`
	ResetPasswordExpires time.Time `bson:"resetPasswordExpires,omitempty"`

	Twitter              vendorOauth `bson:"twitter"`
	Github               vendorOauth `bson:"github"`
	Facebook             vendorOauth `bson:"facebook"`
	Google               vendorOauth `bson:"google"`
	Tumblr               vendorOauth `bson:"tumblr"`
	Search               []string `bson:"search"`
}

func (user *User) CanPlayRoleOf(role string) bool {
	if id_, ok := user.Roles.Account.Id.(bson.ObjectId);
	role == "admin" && ok && len(id_.String()) > 0 {
		return true;
	}
	if id_, ok := user.Roles.Account.Id.(bson.ObjectId);
	role == "account" && ok && len(id_.String()) > 0 {
		return true;
	}
	return false
}

func (user *User) DefaultReturnUrl() (returnUrl string) {
	returnUrl = "/"
	if user.CanPlayRoleOf("account") {
		returnUrl = "/account/"
	}
	if user.CanPlayRoleOf("admin") {
		returnUrl = "/admin/"
	}
	return
}

func (u *User) ValidateUsername(r *Response) {
	u.Username = strings.ToLower(u.Username)
	if len(u.Username) == 0 {
		r.ErrFor["username"] = "required"
	} else {
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\-\_]+$`, u.Username)
		if err != nil {
			println(err.Error())
		}
		if !ok {
			r.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
}

func (u *User) ValidateEmail(r *Response) {
	u.Email = strings.ToLower(u.Email)
	if len(u.Email) == 0 {
		r.ErrFor["email"] = "required"
	} else {
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\+]+@[a-zA-Z0-9\-\_\.]+\.[a-zA-Z0-9\-\_]+$`, u.Email)
		if err != nil {
			println(err.Error())
		}
		if !ok {
			r.ErrFor["email"] = `invalid email format`
		}
	}
}

func (u *User) ValidatePassword(r *Response) {
	if len(u.Password) == 0 {
		r.ErrFor["password"] = "required"
	} else {
		if len(u.Password) < 8 {
			r.ErrFor["password"] = `too weak password`
		}
	}
}

var UserUniqueIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
}

var UserIndex mgo.Index = mgo.Index{
	Key:        []string{"timeCreated", "twitter.id", "github.id", "facebook.id", "google.id", "search"},
}