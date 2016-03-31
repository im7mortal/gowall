package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	User struct{
		   ID bson.ObjectId `bson:"id"`
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
		   ID Status `bson:"id"`
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

func (u *Account) Flow()  {

}



var AccountIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:     "userIndex",
}