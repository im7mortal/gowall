package main

import (
	"gopkg.in/mgo.v2"
)

type Status struct {
	ID string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
	Pivot string `bson:"pivot" json:"pivot"`
}

var StatusIndex mgo.Index = mgo.Index{
	//Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:       "statusIndex",
}