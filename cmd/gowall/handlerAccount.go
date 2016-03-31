package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/im7mortal/gowall/schemas"
	"github.com/gin-gonic/contrib/sessions"
	"encoding/json"
	"net/url"
)

func Account(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func AccountVerification(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func AccountSettingsRender(c *gin.Context) {
	sess := sessions.Default(c)

	public := sess.Get("public")
	session, err := mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	collection := d.C("User")
	us := schemas.User{} // todo pool
	err = collection.FindId(bson.ObjectIdHex(public.(string))).One(&us)
	if err != nil {
		println(err.Error())
	}
	if len(us.Username) != 0 {
		User, _ := json.Marshal(gin.H{
			"_id": us.ID.Hex(),
			"user": us.Username,
			"email": us.Email,
		})
		c.Set("User", url.QueryEscape(string(User)))
	}
	collection = d.C("Account")
	ac := schemas.Account{} // todo pool
	err = collection.FindId(bson.ObjectIdHex(us.Roles.Account.Hex())).One(&ac)
	if err != nil {
		println(err.Error())
	}
	if len(ac.ID) != 0 {
		Account, _ := json.Marshal(gin.H{
			"_id": ac.ID.Hex(),
			"name": gin.H{
				"first": ac.Name.First,
				"middle": ac.Name.Middle,
				"last": ac.Name.Last,
			},
			"company": ac.Company,
			"phone": ac.Phone,
			"zip": ac.Zip,
		})
		c.Set("Account", url.QueryEscape(string(Account)))
	}
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}