package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/contrib/sessions"
	"gopkg.in/mgo.v2"
	"github.com/markbates/goth"
)

func LoginProvider(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("action", "/login/")
	session.Save()
	startOAuth(c)
}

func loginProvider(c *gin.Context, userGoth goth.User) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err := collection.Find(bson.M{userGoth.Provider + ".id": userGoth.UserID}).One(&user)
	// we expect err == mgo.ErrNotFound for success
	if err != nil {
		println("here")
		if err == mgo.ErrNotFound {
			session := sessions.Default(c)
			session.Set("oauthMessage", "No users found linked to your " + userGoth.Provider + " account. You may need to create an account first.")
			session.Save()
			c.Redirect(http.StatusFound, "/login/")
			return
		}
		panic(err)
	}

	session := sessions.Default(c)
	session.Set("public", user.ID.Hex())
	returnURL, ok := session.Get("returnURL").(string)
	session.Delete("returnURL")
	session.Save()

	if ok {
		c.Redirect(http.StatusFound, returnURL)
	} else {
		c.Redirect(http.StatusFound, user.DefaultReturnUrl())
	}
}