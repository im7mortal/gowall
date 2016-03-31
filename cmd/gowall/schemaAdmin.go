package main

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Admin struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	User struct{
		   ID mgo.DBRef `bson:"id"`
		   Name string `bson:"name"`
	   } `bson:"user"`
	Name struct {
		   First string  `bson:"first"`
		   Middle string `bson:"middle"`
		   Last string `bson:"last"`
		   Full string `bson:"full"`
	   } `bson:"name"`
	Groups []AdminGroup `bson:"groups"`
	Permissions []Permission `bson:"permissions"`
	TimeCreated time.Time `bson:"timeCreated"`
	Search []string `bson:"search"`
}

func (u *Admin) Flow()  {

}

func (admin *Admin) HasPermissionTo(requiredPermission string) (hasPermission bool) {
	hasPermission = false
	//check group permissions
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

func (u *Admin) IsMemberOf(groupName string) bool {
	for _, group := range u.Groups{
		if group.Name == groupName {
			return true
		}
	}
	return false
}

var AdminIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:     "userIndex",
}