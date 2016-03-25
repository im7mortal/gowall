package schemas

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Schema interface{
	Flow()
}

type Account struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	User struct{
			ID mgo.DBRef `bson:"id"`
			Name string `bson:"name"`
		} `bson:"user"`
	IsVerified string `bson:"isVerified"`
	VerificationToken string `bson:"verificationToken"`
	Name struct {
			First string  `bson:"first"`
			Middle string `bson:"middle"`
			Last string `bson:"last"`
			Full string `bson:"full"`
		} `bson:"name"`
	Company string `bson:"company"`
	Phone string `bson:"phone"`
	Zip string `bson:"zip"`
	Status struct {
		ID mgo.DBRef `bson:"id"`
		Name string `bson:"name"`
		UserCreated struct {
			ID mgo.DBRef `bson:"id"`
			Name string `bson:"name"`
			Time time.Time `bson:"time"`
		} `bson:"userCreated"`
	} `bson:"status"`
	StatusLog []StatusLog `bson:"statusLog"`
	Notes []Note `bson:"notes"`
	UserCreated struct{
			ID mgo.DBRef `bson:"id"`
			Name string `bson:"name"`
			Time time.Time `bson:"time"`
		} `bson:"userCreated"`
	Search []string `bson:"search"`
}

type StatusLog struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	ID_ mgo.DBRef `bson:"id"`
	Name string `bson:"name"`
	UserCreated struct{
		ID mgo.DBRef `bson:"id"`
		Name string `bson:"name"`
		Time string `bson:"time"`
	} `bson:"userCreated"`
}

type Note struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Data string `bson:"data"`
	UserCreated struct{
		ID mgo.DBRef `bson:"id"`
		Name string `bson:"name"`
		Time string `bson:"time"`
	} `bson:"userCreated"`
}

type Status struct {
	ID    bson.ObjectId `bson:"_id"`
	Pivot string `bson:"pivot"`
	Name  string `bson:"name"`
}

type Status1 struct {
	ID    bson.ObjectId `bson:"_id"`
	Pivot string `bson:"pivot"`
	Name  string `bson:"name"`
}



type vendorOauth struct {

}

type User struct {// todo uniq
	ID                   bson.ObjectId `bson:"_id"`
	Username             string `bson:"username"`
	Password             string `bson:"password"`
	Email                string `bson:"email"`
	Roles                struct {
						 Admin   mgo.DBRef `bson:"name,omitempty"`
						 Account mgo.DBRef `bson:"time,omitempty"`
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



var UserIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:     "userIndex",
}