package schemas

import (
	"gopkg.in/mgo.v2"
)

type Status struct {
	ID    string `bson:"_id"`
	Pivot string `bson:"pivot"`
	Name  string `bson:"name"`
}

var StatusIndex mgo.Index = mgo.Index{
	//Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:       "statusIndex",
}