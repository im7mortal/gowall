package main

import (
	"gopkg.in/mgo.v2"
)

type StatusLog struct {
	ID mgo.DBRef `bson:"id"`
	Name string `bson:"name"`
	UserCreated struct{
		   ID mgo.DBRef `bson:"id"`
		   Name string `bson:"name"`
		   Time string `bson:"time"`
	   } `bson:"userCreated"`
}
