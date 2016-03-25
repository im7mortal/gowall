package schemas

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
)

/*
{
    user: {
      id: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
      name: { type: String, default: '' }
    },
    name: {
      full: { type: String, default: '' },
      first: { type: String, default: '' },
      middle: { type: String, default: '' },
      last: { type: String, default: '' },
    },
    groups: [{ type: String, ref: 'AdminGroup' }],
    permissions: [{
      name: String,
      permit: Boolean
    }],
    timeCreated: { type: Date, default: Date.now },
    search: [String]
  }
*/

type Permission struct {
	Name string `bson:"name"`
	Permit bool `bson:"permit"`
}

type Admin struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	User struct{
		   ID mgo.DBRef `bson:"id"`
		   Name string `bson:"name"`
	   } `bson:"user"`
	Name struct {
		   First string  `bson:"first"`
		   Middle string `bson:"middle"`
		   Last string `bson:"last"`
		   Full string `bson:"full"`
	   } `bson:"name"`
	Groups []mgo.DBRef `bson:"groups"`

	Permissions []Permission `bson:"permissions"`
	TimeCreated time.Time `bson:"timeCreated"`
	Search []string `bson:"search"`
}

func (u *Admin) Flow()  {

}

var UserIndex mgo.Index = mgo.Index{
	Key:        []string{"username", "email"},
	Unique:     true,
	DropDups:   true,
	Background: true,
	Sparse:     true,
	Name:     "userIndex",
}