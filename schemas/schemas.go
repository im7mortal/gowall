package main

import (
	"gopkg.in/mgo.v2"
	"time"
)


/*
{
    user: {
      id: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
      name: { type: String, default: '' }
    },
    isVerified: { type: String, default: '' },
    verificationToken: { type: String, default: '' },
    name: {
      first: { type: String, default: '' },
      middle: { type: String, default: '' },
      last: { type: String, default: '' },
      full: { type: String, default: '' }
    },
    company: { type: String, default: '' },
    phone: { type: String, default: '' },
    zip: { type: String, default: '' },
    status: {
      id: { type: String, ref: 'Status' }, //todo
      name: { type: String, default: '' },
      userCreated: {
        id: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
        name: { type: String, default: '' },
        time: { type: Date, default: Date.now }
      }
    },
    statusLog: [mongoose.modelSchemas.StatusLog],
    notes: [mongoose.modelSchemas.Note],
    userCreated: {
      id: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
      name: { type: String, default: '' },
      time: { type: Date, default: Date.now }
    },
    search: [String]
  }

*/

type Account struct {
	User struct{
			ID mgo.DBRef `bson:"id"`
			Name string `bson:"name"`
		} `bson:"user"`
	IsVerified string `bson:"isVerified"`
	VerificationToken string `bson:"verificationToken"`
	Name struct {
			First string  `bson:"first"`
			Middle string `bson:"middle"`
			Last string `bson:"last"`
			Full string `bson:"full"`
		} `bson:"name"`
	Company string `bson:"company"`
	Phone string `bson:"phone"`
	Zip string `bson:"zip"`
	Status struct {
		ID mgo.DBRef `bson:"id"`
		Name string `bson:"name"`
		UserCreated struct {
			ID mgo.DBRef `bson:"id"`
			Name string `bson:"name"`
			Time time.Time `bson:"time"`
		} `bson:"userCreated"`
	} `bson:"status"`
	StatusLog []StatusLog `bson:"statusLog"`
	Notes []Note `bson:"notes"`
	UserCreated struct{
			ID mgo.DBRef `bson:"id"`
			Name string `bson:"name"`
			Time time.Time `bson:"time"`
		} `bson:"userCreated"`
	Search []string `bson:"search"`
}

type StatusLog struct {
	ID mgo.DBRef `bson:"id"`
	Name string `bson:"name"`
	UserCreated struct{
		ID mgo.DBRef `bson:"id"`
		Name string `bson:"name"`
		Time string `bson:"time"`
	} `bson:"userCreated"`
	Search []string `bson:"search"`
}

type Note struct {
	Data string `bson:"data"`
	UserCreated struct{
		ID mgo.DBRef `bson:"id"`
		Name string `bson:"name"`
		Time string `bson:"time"`
	} `bson:"userCreated"`
	Search []string `bson:"search"`
}


