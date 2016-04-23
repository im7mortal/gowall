package main

import "gopkg.in/mgo.v2/bson"

type Category struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Pivot string `bson:"pivot"`
}
