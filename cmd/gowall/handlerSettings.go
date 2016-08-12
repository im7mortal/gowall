package main

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

func renderAccountSettings(c *gin.Context) {
	sess := sessions.Default(c)

	user := getUser(c)
	injectSocials(c)
	doUserHasSocials(c, user)

	public := sess.Get("public")
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user = &User{}
	err := collection.FindId(bson.ObjectIdHex(public.(string))).One(user)
	if err != nil {
		EXCEPTION(err)
	}
	if len(user.Username) != 0 {
		User, _ := json.Marshal(gin.H{
			"_id":      user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
		})
		c.Set("User", template.JS(getEscapedString(string(User))))
	}
	collection = db.C(ACCOUNTS)
	ac := Account{}
	err = collection.FindId(user.Roles.Account).One(&ac)
	if err != nil {
		EXCEPTION(err)
	}
	if len(ac.ID) != 0 {
		Account, _ := json.Marshal(gin.H{
			"_id": ac.ID.Hex(),
			"name": gin.H{
				"first":  ac.Name.First,
				"middle": ac.Name.Middle,
				"last":   ac.Name.Last,
			},
			"company": ac.Company,
			"phone":   ac.Phone,
			"zip":     ac.Zip,
		})
		c.Set("Account", template.JS(getEscapedString(string(Account))))
	}
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func setSettings(c *gin.Context) {
	account := getAccount(c)
	response := newResponse(c)
	err := account.changeData(response)
	if err != nil {
		response.Fail()
		return
	}
	response.Finish()
}

func changePassword(c *gin.Context) {
	user := getUser(c)
	response := newResponse(c)
	err := user.changePassword(response)
	if err != nil {
		response.Fail()
		return
	}
	response.Finish()
}

func changeIdentity(c *gin.Context) {
	user := getUser(c)
	response := newResponse(c)
	err := user.changeIdentity(response)
	if err != nil {
		response.Fail()
		return
	}
	response.Finish()
}

func settingsProvider_(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("action", "/account/settings/")
	session.Save()
	startOAuth(c)
}

func settingsProvider(c *gin.Context, userGoth goth.User) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := &User{}
	err := collection.Find(bson.M{userGoth.Provider + ".id": userGoth.UserID}).One(user)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		session := sessions.Default(c)
		session.Set("oauthMessage", "Another user has already connected with that " + userGoth.Provider + " account")
		session.Save()
		c.Redirect(http.StatusFound, "/account/settings/")
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	user = getUser(c)

	user.updateProvider(userGoth)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		EXCEPTION(err)
	}

	c.Redirect(http.StatusFound, "/account/settings/")
}

func disconnectProvider(c *gin.Context) {
	user := getUser(c)
	user.disconnectProviderDB(c.Param("provider"))
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err := collection.UpdateId(user.ID, user)
	if err != nil {
		EXCEPTION(err)
	}

	c.Redirect(http.StatusFound, "/account/settings/")
}
