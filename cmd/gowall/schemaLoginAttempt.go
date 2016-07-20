package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func init() {
	expiration, err := time.ParseDuration(config.LoginAttempts.LogExpiration)
	if err != nil {
		panic(err.Error()) // todo FATAL
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
	Key:         []string{"username", "email"},
	Unique:      true,
	ExpireAfter: LoginAttemptsExpireAfter,
}
