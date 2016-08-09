package main

import (
	"gopkg.in/mgo.v2/bson"
)

type accountStatus struct {
	ID          bson.ObjectId `json:"id" bson:"id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	UserCreated struct {
					ID   bson.ObjectId `json:"id" bson:"id,omitempty"`
					Name string        `json:"name" bson:"name"`
					Time string        `json:"time" bson:"time"`
				} `json:"userCreated" bson:"userCreated"`
}
