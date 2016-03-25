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

type Admin struct {// todo uniq
	ID                   bson.ObjectId `bson:"_id"`
	Username             string `bson:"username"`
	Password             string `bson:"password"`
	Email                string `bson:"email"`
	Roles                struct {
						 Admin   mgo.DBRef `bson:"name,omitempty"`
						 Account mgo.DBRef `bson:"time,omitempty"`
					 } `bson:"roles"`

	IsActive             string `bson:"isActive,omitempty"`
	TimeCreated          time.Time `bson:"timeCreated"`
	ResetPasswordToken   string `bson:"resetPasswordToken,omitempty"`
	ResetPasswordExpires time.Time `bson:"resetPasswordExpires,omitempty"`

	Twitter              vendorOauth `bson:"twitter"`
	Github               vendorOauth `bson:"github"`
	Facebook             vendorOauth `bson:"facebook"`
	Google               vendorOauth `bson:"google"`
	Tumblr               vendorOauth `bson:"tumblr"`
	Search               []string `bson:"search"`
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