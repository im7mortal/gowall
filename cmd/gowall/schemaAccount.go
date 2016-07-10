package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	User struct{
		   ID bson.ObjectId `bson:"id" json:"id"`
		   Name string `bson:"name" json:"name"`
	   } `bson:"user" json:"user"`
	IsVerified string `bson:"isVerified" json:"isVerified"`
	VerificationToken string `bson:"verificationToken" json:"verificationToken"`
	Name struct {
		   First string  `bson:"first" json:"first"`
		   Middle string `bson:"middle" json:"middle"`
		   Last string `bson:"last" json:"last"`
		   Full string `bson:"full" json:"full"`
	   } `bson:"name" json:"name"`
	Company string `bson:"company" json:"company"`
	Phone string `bson:"phone" json:"phone"`
	Zip string `bson:"zip" json:"zip"`
	Status struct {
		   ID Status `bson:"id" json:"id"`
		   Name string `bson:"name" json:"name"`
		   UserCreated struct {
				 ID mgo.DBRef `bson:"id" json:"id"`
				 Name string `bson:"name" json:"name"`
				 Time time.Time `bson:"time" json:"time"`
			 } `bson:"userCreated" json:"userCreated"`
	   } `bson:"status" json:"status"`
	StatusLog []StatusLog `bson:"statusLog" json:"statusLog"`
	Notes []Note `bson:"notes" json:"notes"`
	UserCreated struct{
		   ID mgo.DBRef `bson:"id" json:"id"`
		   Name string `bson:"name" json:"name"`
		   Time time.Time `bson:"time" json:"time"`
	   } `bson:"userCreated" json:"userCreated"`
	Search []string `bson:"search" json:"search"`
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