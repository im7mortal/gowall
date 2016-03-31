package main

import "gopkg.in/mgo.v2/bson"

type LoginAttempt struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	IP string `bson:"ip"`
	User string `bson:"user"`
/*	Time struct{
		   ID bson.ObjectId `bson:"id"`
		   Name string `bson:"name"`
		   Time string `bson:"time"`
	   } `bson:"time"`*/
}