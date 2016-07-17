package main

import (
	"gopkg.in/mgo.v2/bson"
)

type StatusLog struct {
	ID bson.ObjectId `bson:"id"`
	Name string `bson:"name"`
	UserCreated struct{
		   ID bson.ObjectId `bson:"id"`
		   Name string `bson:"name"`
		   Time string `bson:"time"`
	   } `bson:"userCreated"`
}
