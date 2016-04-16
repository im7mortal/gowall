package main

import (
	"gopkg.in/mgo.v2"
	"strings"
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

// attempt to get dbName from URL
// it will work on MongoLab where dbName is part of url
func getDBName(url *string) string {
	arr := strings.Split(*url, ":")
	arr = strings.Split(arr[len(arr) - 1], "/")
	return arr[len(arr) - 1]
}