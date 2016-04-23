package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Permission struct {
	Name string `bson:"name"`
	Permit bool `bson:"permit"`
}

type AdminGroup struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string `bson:"name"`
	Permissions []Permission `bson:"permissions"`
}

func (u *AdminGroup) Flow()  {

}



var AdminGroupIndex mgo.Index = mgo.Index{
	Key:        []string{"name"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:       "adminGroupIndex",
}