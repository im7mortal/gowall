package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Account struct {
	ID                bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	User              struct {
						  ID   bson.ObjectId `bson:"id,omitempty" json:"id"`
						  Name string        `bson:"name" json:"name"`
					  } `bson:"user" json:"user"`
	IsVerified        string `bson:"isVerified" json:"isVerified"`
	VerificationToken string `bson:"verificationToken" json:"-"`
	Name              struct {
						  First  string `bson:"first" json:"first"`
						  Middle string `bson:"middle" json:"middle"`
						  Last   string `bson:"last" json:"last"`
						  Full   string `bson:"full" json:"full"`
					  } `bson:"name" json:"name"`
	Company           string          `bson:"company" json:"company"`
	Phone             string          `bson:"phone" json:"phone"`
	Zip               string          `bson:"zip" json:"zip"`
	Status            accountStatus   `bson:"status" json:"status"`
	StatusLog         []accountStatus `bson:"statusLog" json:"statusLog"`
	Notes             []Note          `bson:"notes" json:"notes"`
	UserCreated       struct {
						  ID   bson.ObjectId `bson:"id,omitempty" json:"id"`
						  Name string        `bson:"name" json:"name"`
						  Time time.Time     `bson:"time" json:"time"`
					  } `bson:"userCreated" json:"userCreated"`
	Search            []string `bson:"search" json:"search"`
}

func (a *Account) DecodeRequest(c *gin.Context) {
	err := json.NewDecoder(c.Request.Body).Decode(a)
	if err != nil {
		EXCEPTION(err)
	}
	return
}

var AccountIndex mgo.Index = mgo.Index{
	Key: []string{"user", "status.id", "search"},
}

func (a *Account) linkUser(db *mgo.Database, user *User) (err error) {

	// patchUser
	collection := db.C(USERS)
	err = collection.UpdateId(user.ID, bson.M{
		"$set": bson.M{"roles.account": a.ID},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}

	// patchAccount
	collection = db.C(ACCOUNTS)
	err = collection.UpdateId(a.ID, bson.M{
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

func (a *Account) unlinkUser(db *mgo.Database, user *User) (err error) {

	// patchUser
	collection := db.C(USERS)
	err = collection.Update(bson.M{"roles.account": user.Roles.Account}, bson.M{
		"$set": bson.M{"roles.account": ""},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}

	// patchAccount
	collection = db.C(ACCOUNTS)
	err = collection.UpdateId(user.Roles.Account, bson.M{
		"$set": bson.M{"roles.account": ""},
	})

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		return
	}
	return
}


func (a *Account) changeData(r *Response) (err error) {

	var body struct {
		First   string `json:"first"`
		Middle  string `json:"middle"`
		Last    string `json:"last"`
		Company string `json:"company"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
	}
	err = json.NewDecoder(r.c.Request.Body).Decode(&body)

	if len(body.First) == 0 {
		r.ErrFor["first"] = "required"
	}
	if len(body.Last) == 0 {
		r.ErrFor["last"] = "required"
	}

	if r.HasErrors() {
		err = Err
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()

	a.Name.Full = body.First + " " + body.Last
	a.Name.First = body.First
	a.Name.Middle = body.Middle
	a.Name.Last = body.Last
	a.Company = body.Company
	a.Phone = body.Phone
	a.Zip = body.Zip
	a.Search = a.Search[:0]
	a.Search = append(a.Search,
		body.First,
		body.Middle,
		body.Last,
		body.Company,
		body.Phone,
		body.Zip,
	)

	collection := db.C(ACCOUNTS)
	err = collection.UpdateId(a.ID, a)
	if err != nil {
		EXCEPTION(err)
	}

	return
}
