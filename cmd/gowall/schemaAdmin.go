package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type Admin struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	User struct{
		   ID bson.ObjectId `bson:"id,omitempty" json:"id"`
		   Name string `bson:"name" json:"name"`
	   } `bson:"user" json:"user"`
	Name struct {
		   First string  `bson:"first" json:"first"`
		   Middle string `bson:"middle" json:"middle"`
		   Last string `bson:"last" json:"last"`
		   Full string `bson:"full" json:"full"`
	   } `bson:"name" json:"name"`
	Groups []string `bson:"groups" json:"groups"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
	TimeCreated time.Time `bson:"timeCreated" json:"timeCreated"`
	Search []string `bson:"search" json:"search"`
}

func (u *Admin) Flow()  {

}

func (a *Admin) DecodeRequest(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(a)
	if err != nil {
		panic(err)
	}
	return
}

func (admin *Admin) HasPermissionTo(requiredPermission string) (hasPermission bool) {
	hasPermission = false
	//check group permissions
	/*
	for _, adminGroup := range admin.Groups{
		for _, permission := range adminGroup.Permissions{
			if permission.Name == requiredPermission {
				hasPermission = true
				break
			}
		}
		if hasPermission {
			break
		}
	}
	*/

	//check admin permissions
	for _, permission := range admin.Permissions{
		if permission.Name == requiredPermission {
			if permission.Permit {
				return true
			}
			return false
		}
	}
	return
}

func (admin *Admin) IsMemberOf(groupName string) bool {
	for _, group := range admin.Groups{
		if group == groupName {
			return true
		}
	}
	return false
}

var AdminsIndex mgo.Index = mgo.Index{
	Key:        []string{"user.id", "search"},
}