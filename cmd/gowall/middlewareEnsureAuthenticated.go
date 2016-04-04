package main

import (
"github.com/gin-gonic/contrib/sessions"
"github.com/gin-gonic/gin"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
	"net/http"
)

const MONGOURL  = "mongodb://localhost:27017"
const DBNAME  = "test"
const USERTABLE  = "User"

func EnsureAuthenticated(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		c.Next()
	} else {
		session := sessions.Default(c)
		session.Set("returnURL", c.Request.URL.Path)
		session.Save()
		c.Redirect(http.StatusFound, "/login/")
	}
}

func IsAuthenticated(c *gin.Context) {
	c.Set("isAuthenticated", false)

	session := sessions.Default(c)

	public := session.Get("public")

	if public != nil && len(public.(string)) > 0 {
		session, err := mgo.Dial(MONGOURL)
		defer session.Close()
		if err != nil {
			println(err.Error())
		}
		collection := session.DB(DBNAME).C(USERTABLE)
		us := UsersPool.Get().(*User)
		defer UsersPool.Put(us)
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(public.(string))}).One(&us)
		if err != nil {
			println(err.Error())
		}
		if len(us.Username) > 0 {
			c.Set("Logined", true) // todo what is different between "Logined" and "isAuthenticated"
			c.Set("isAuthenticated", true)
			c.Set("UserName", us.Username)
			c.Set("CurUser", us)
			c.Set("DefaultReturnUrl", us.DefaultReturnUrl()) // todo
		}
	}

	c.Next()
}
