package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Admin struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	User        struct {
					ID   bson.ObjectId `bson:"id,omitempty" json:"id"`
					Name string        `bson:"name" json:"name"`
				} `bson:"user" json:"user"`
	Name        struct {
					First  string `bson:"first" json:"first"`
					Middle string `bson:"middle" json:"middle"`
					Last   string `bson:"last" json:"last"`
					Full   string `bson:"full" json:"full"`
				} `bson:"name" json:"name"`
	Groups      []string     `bson:"groups" json:"-"` // drywall author used mutable objects
	GroupsJS    []AdminGroup `json:"groups" bson:"-"` // I separate types for mongoDB and json
	Permissions []Permission `bson:"permissions" json:"permissions"`
	TimeCreated time.Time    `bson:"timeCreated" json:"timeCreated"`
	Search      []string     `bson:"search" json:"search"`
}

func (u *Admin) Flow() {

}

func (a *Admin) DecodeRequest(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(a)
	if err != nil {
		EXCEPTION(err)
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
	for _, permission := range admin.Permissions {
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
	for _, group := range admin.Groups {
		if group == groupName {
			return true
		}
	}
	return false
}

var AdminsIndex mgo.Index = mgo.Index{
	Key: []string{"user.id", "search"},
}

func (admin *Admin) linkUser(db *mgo.Database, user *User) (err error) {

	// patchUser
	collection := db.C(USERS)
	err = collection.UpdateId(user.ID, bson.M{
		"$set": bson.M{"roles.admin": admin.ID},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}

	// patchAdministrator
	collection = db.C(ADMINS)
	err = collection.UpdateId(admin.ID, bson.M{
		"$set": bson.M{"user": bson.M{
			"id":   user.ID,
			"name": user.Username,
		}},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}
	return
}

func (admin *Admin) unlinkUser(db *mgo.Database, user *User) (err error) {

	// patchUser
	collection := db.C(USERS)
	err = collection.Update(bson.M{"roles.admin": user.Roles.Admin}, bson.M{
		"$set": bson.M{"roles.admin": ""},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}

	// patchAdministrator
	collection = db.C(ADMINS)
	err = collection.UpdateId(user.Roles.Admin, bson.M{
		"$set": bson.M{"user": bson.M{}},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}
	return
}

func (admin *Admin) populateGroups(db *mgo.Database) (adminGroups []AdminGroup, err error) {
	// populateGroups
	collection := db.C(ADMINGROUPS)
	adminGroups = []AdminGroup{}
	err = collection.Find(nil).All(&adminGroups)
	if err != nil {
		return
	}
	admin.GroupsJS = []AdminGroup{}
	for _, adminGroupID := range admin.Groups {
		for _, adminGroup := range adminGroups {
			if adminGroupID == adminGroup.ID {
				admin.GroupsJS = append(admin.GroupsJS, adminGroup)
			}
		}
	}
	return
}
