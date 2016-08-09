package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Permission struct {
	Name   string        `bson:"name" json:"name"`
	Permit bool          `bson:"permit" json:"permit"`
	ID     bson.ObjectId `bson:"_id,omitempty" json:"_id"`
}

type AdminGroup struct {
	ID          string       `bson:"_id,omitempty" json:"_id"`
	Name        string       `bson:"name" json:"name"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
}

func (a *AdminGroup) DecodeRequest(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(a)
	if err != nil {
		EXCEPTION(err)
	}
	return
}

var AdminGroupIndex mgo.Index = mgo.Index{
	Key:    []string{"name"},
	Unique: true,
}
