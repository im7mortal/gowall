package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Status struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string `bson:"name"`
	Pivot string `bson:"pivot"`
}

var StatusIndex mgo.Index = mgo.Index{
	//Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:       "statusIndex",
}