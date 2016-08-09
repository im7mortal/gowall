package main

import "gopkg.in/mgo.v2"

type Category struct {
	ID    string `bson:"_id" json:"_id"`
	Name  string `bson:"name" json:"name"`
	Pivot string `bson:"pivot" json:"pivot"`
}

var CategoryIndex mgo.Index = mgo.Index{
	Key:    []string{"name", "pivot"},
}
