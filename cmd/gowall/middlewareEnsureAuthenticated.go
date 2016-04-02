package main

import (
"github.com/gin-gonic/contrib/sessions"
"github.com/gin-gonic/gin"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
	"net/http"
)

func IsAuthenticated(c *gin.Context) {
	isAuthenticated := false
	sess := sessions.Default(c)

	public := sess.Get("public")

	if public != nil && len(public.(string)) > 0 {
		session, err := mgo.Dial("mongodb://localhost:27017")
		defer session.Close()
		if err != nil {
			println(err.Error())
		}
		d := session.DB("test")
		collection := d.C("User")
		us := User{}
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(public.(string))}).One(&us)
		if err != nil {
			println(err.Error())
		}
		if len(us.Username) > 0 {
			isAuthenticated = true
			c.Set("Logined", true) // todo what is different between "Logined" and "isAuthenticated"
			c.Set("isAuthenticated", true)
			c.Set("UserName", us.Username)
			c.Set("CurUser", us)
			c.Set("DefaultReturnUrl", us.DefaultReturnUrl()) // todo
		}
	}
	if !isAuthenticated {
		goToLogin := false
		if len(c.Request.URL.Path) >= 7 && c.Request.URL.Path[:7] == "/admin/" {
			goToLogin = true
		}
		if len(c.Request.URL.Path) >= 9 && c.Request.URL.Path[:9] == "/account/" {
			goToLogin = true
		}
		if goToLogin {
			session := sessions.Default(c)
			session.Set("returnURL", c.Request.URL.Path)
			session.Save()
			c.Redirect(http.StatusFound, "/login/")
			return
		}
	}
	c.Next()
}

func EnsureAuthenticated(c *gin.Context) {
	isAuthenticated := false

	session := sessions.Default(c)

	public := string(session.Get("public"))

	if public != nil && len(public.(string)) > 0 {
		session, err := mgo.Dial("mongodb://localhost:27017")
		defer session.Close()
		if err != nil {
			println(err.Error())
		}
		d := session.DB("test")
		collection := d.C("User")
		us := User{}
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(public.(string))}).One(&us)
		if err != nil {
			println(err.Error())
		}
		if len(us.Username) > 0 {
			isAuthenticated = true
			c.Set("Logined", true) // todo what is different between "Logined" and "isAuthenticated"
			c.Set("isAuthenticated", true)
			c.Set("UserName", us.Username)
			c.Set("CurUser", us)
			c.Set("DefaultReturnUrl", us.DefaultReturnUrl()) // todo
		}
	}
	if !isAuthenticated {
		goToLogin := false
		if len(c.Request.URL.Path) >= 7 && c.Request.URL.Path[:7] == "/admin/" {
			goToLogin = true
		}
		if len(c.Request.URL.Path) >= 9 && c.Request.URL.Path[:9] == "/account/" {
			goToLogin = true
		}
		if goToLogin {
			session := sessions.Default(c)
			session.Set("returnURL", c.Request.URL.Path)
			session.Save()
			c.Redirect(http.StatusFound, "/login/")
			return
		}
	}
	c.Next()
}
