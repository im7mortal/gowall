package main

import (
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2/bson"
)

const urlMONGO  = "mongodb://localhost:27017"



/*
func init() {

	session, err := mgo.Dial(urlMONGO)
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



	err = c.EnsureIndex(UserIndex)
	if err != nil {
		println(err.Error())
	}

	err = c.Insert(User{Username:"valera", ID: bson.NewObjectId()})
	if err != nil {
		println(err.Error())
	}
	us := User{}
	_ = c.Find(bson.M{"username": "valera"}).One(&us)

	println(us.ID.Hex())

	acc := Account{ID: bson.NewObjectId()}
	*/
/*
	acc.User.ID = mgo.DBRef{
		Id: us.ID,
		Collection: c.Name,
		Database: d.Name,
	}
	*//*

	err = c.Insert(acc)
	if err != nil {
		println(err.Error())
	}



	println("dddddddddddddddddddddddddddddddddddddddddd")
	println(i)
	println("dddddddddddddddddddddddddddddddddddddddddd")

}*/
