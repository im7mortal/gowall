package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Note struct {
	ID          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Data        string        `json:"data" bson:"data"`
	UserCreated struct {
					ID   bson.ObjectId `json:"id" bson:"id,omitempty"`
					Name string        `json:"name" bson:"name"`
					Time time.Time     `json:"time" bson:"time"`
				} `json:"userCreated" bson:"userCreated"`
}
