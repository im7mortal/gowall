package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type vendorOauth struct {

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
						 Admin   bson.ObjectId `bson:"name,omitempty"`
						 Account bson.ObjectId `bson:"time,omitempty"`
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
	if  role == "account" && len(user.Roles.Account) > 0 {
	return true;
	}
	return false
}

func (user *User) DefaultReturnUrl() (returnUrl string) {
	returnUrl = "/"
	if user.CanPlayRoleOf("admin") {
		returnUrl = "/admin/"
	}
	if user.CanPlayRoleOf("account") {
		returnUrl = "/account/"
	}
	return
}

/*func (user *User) EncryptPassword(password string, done bool) (err error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {

	}

}*/



var UserIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:     "userIndex",
}