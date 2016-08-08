package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
)

func init() {
	expiration, err := time.ParseDuration(config.LoginAttempts.LogExpiration)
	if err != nil {
		log.Fatal("LoginAttemptsExpireAfter is not set: ", err)
	}
	LoginAttemptsExpireAfter = expiration
}

var LoginAttemptsExpireAfter time.Duration

type LoginAttempt struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	IP   string        `bson:"ip"`
	User string        `bson:"user"`
	Time time.Time     `bson:"time"`
}

var LoginAttemptsIndex mgo.Index = mgo.Index{
	Key:         []string{"ip", "user"},
	ExpireAfter: LoginAttemptsExpireAfter,
}
