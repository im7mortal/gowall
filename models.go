package main

import (
	"gopkg.in/mgo.v2"
	"github.com/im7mortal/gowall/schemas"
	//"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/bson"
)

const url  = "mongodb://localhost:27017"



func init() {

	session, err := mgo.Dial(url)
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	c := d.C("restaurants")
	//c.DropCollection()
	c.Create(&mgo.CollectionInfo{})
	i, err := c.Count()

	if err != nil {
		println(err.Error())
	}



	err = c.EnsureIndex(schemas.UserIndex)
	if err != nil {
		println(err.Error())
	}

	err = c.Insert(schemas.User{Username:"valera", ID: bson.NewObjectId()})
	if err != nil {
		println(err.Error())
	}
	us := schemas.User{}
	_ = c.Find(bson.M{"username": "valera"}).One(&us)

	println(us.ID.Hex())

	acc := schemas.Account{ID: bson.NewObjectId()}
	/*
	acc.User.ID = mgo.DBRef{
		Id: us.ID,
		Collection: c.Name,
		Database: d.Name,
	}
	*/
	err = c.Insert(acc)
	if err != nil {
		println(err.Error())
	}



	println("dddddddddddddddddddddddddddddddddddddddddd")
	println(i)
	println("dddddddddddddddddddddddddddddddddddddddddd")

}