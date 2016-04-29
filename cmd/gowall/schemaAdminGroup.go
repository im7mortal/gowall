package main

import (
	"gopkg.in/mgo.v2"
)

type Permission struct {
	Name string `bson:"name" json:"name"`
	Permit bool `bson:"permit" json:"permit"`
}

type AdminGroup struct {
	ID          string `bson:"_id,omitempty" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
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