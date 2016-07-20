package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Account struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	User struct {
		ID   bson.ObjectId `bson:"id" json:"id"`
		Name string        `bson:"name" json:"name"`
	} `bson:"user" json:"user"`
	IsVerified        string `bson:"isVerified" json:"isVerified"`
	VerificationToken string `bson:"verificationToken" json:"verificationToken"`
	Name              struct {
		First  string `bson:"first" json:"first"`
		Middle string `bson:"middle" json:"middle"`
		Last   string `bson:"last" json:"last"`
		Full   string `bson:"full" json:"full"`
	} `bson:"name" json:"name"`
	Company string `bson:"company" json:"company"`
	Phone   string `bson:"phone" json:"phone"`
	Zip     string `bson:"zip" json:"zip"`
	Status  struct {
		ID          bson.ObjectId `bson:"id" json:"id"`
		Name        string        `bson:"name" json:"name"`
		UserCreated struct {
			ID   bson.ObjectId `bson:"id" json:"id"`
			Name string        `bson:"name" json:"name"`
			Time time.Time     `bson:"time" json:"time"`
		} `bson:"userCreated" json:"userCreated"`
	} `bson:"status" json:"status"`
	StatusLog   []StatusLog `bson:"statusLog" json:"statusLog"`
	Notes       []Note      `bson:"notes" json:"notes"`
	UserCreated struct {
		ID   bson.ObjectId `bson:"id" json:"id"`
		Name string        `bson:"name" json:"name"`
		Time time.Time     `bson:"time" json:"time"`
	} `bson:"userCreated" json:"userCreated"`
	Search []string `bson:"search" json:"search"`
}

func (a *Account) DecodeRequest(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(a)
	if err != nil {
		panic(err)
	}
	return
}

var AccountIndex mgo.Index = mgo.Index{
	Key: []string{"user", "status.id", "search"},
}

func (a *Account) linkUser(db *mgo.Database, user User) (err error) {

	// patchUser
	collection := db.C(USERS)
	err = collection.UpdateId(user.ID, bson.M{
		"$set": bson.M{"roles.admin": a.ID},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		return
	}

	// patchAdministrator
	collection = db.C(ACCOUNTS)
	err = collection.UpdateId(a.ID, bson.M{
		"$set": bson.M{"user": bson.M{
			"id":   user.ID,
			"name": user.Username,
		}},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		return
	}
	return
}

func (a *Account) unlinkUser(db *mgo.Database, user User) (err error) {

	// patchUser
	collection := db.C(USERS)
	err = collection.Update(bson.M{"roles.admin": a.ID}, bson.M{
		"$set": bson.M{"roles.admin": ""},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		return
	}

	// patchAdministrator
	collection = db.C(ADMINS)
	err = collection.UpdateId(a.ID, bson.M{
		"$set": bson.M{"user": bson.M{}},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		return
	}
	return
}
