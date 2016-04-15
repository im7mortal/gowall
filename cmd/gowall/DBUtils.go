package main

import (
	"gopkg.in/mgo.v2"
)

func getMongoDBInstance() *mgo.Database {
	session, err := mgo.Dial(config.MongoDB)
	if err != nil {
		panic(err)
	}
	// if MongoDBName == "" it will check the connection url MongoDB for a dbname
	// that logic inside mgo
	return session.DB(config.dbName)
}
