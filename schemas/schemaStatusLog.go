package schemas

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type StatusLog struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	ID_ mgo.DBRef `bson:"id"`
	Name string `bson:"name"`
	UserCreated struct{
		   ID mgo.DBRef `bson:"id"`
		   Name string `bson:"name"`
		   Time string `bson:"time"`
	   } `bson:"userCreated"`
}
