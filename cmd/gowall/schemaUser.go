package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"strings"
	"regexp"
)

type vendorOauth struct {
	ID string `bson:"id"`
}

var UsersPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

type User struct {// todo uniq
	ID                   bson.ObjectId `bson:"_id"`
	Username             string `bson:"username"`
	Password             string `bson:"password"`
	Email                string `bson:"email"`
	Roles                struct {
						 Admin   bson.ObjectId `bson:"admin,omitempty"`
						 Account bson.ObjectId `bson:"account,omitempty"`
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

func (u *User) Flow()  {

}

func (user *User) CanPlayRoleOf(role string) bool {
	if role == "admin" && len(user.Roles.Admin) > 0 {
		return true;
	}
	if role == "account" && len(user.Roles.Account) > 0 {
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

/*func (user *User) EncryptPassword(password string, done bool) (err error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {

	}

}*/

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

var UserIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:     "userIndex",
}