package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/contrib/sessions"
	"gopkg.in/mgo.v2"
)

func LoginRender(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		var redirectURL string
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		redirectURL = defaultReturnUrl.(string)
		session := sessions.Default(c)
		returnURL := session.Get("returnURL")
		if returnURL != nil {
			redirectURL = returnURL.(string)
			session.Delete("returnURL")
			session.Save()
		}
		c.Redirect(http.StatusFound, redirectURL)
	} else {
		render, _ := TemplateStorage[c.Request.URL.Path]

		_, oauthTwitter := config.Socials["twitter"]
		_, oauthGitHub := config.Socials["github"]
		_, oauthFacebook := config.Socials["facebook"]
		_, oauthGoogle := config.Socials["google"]
		_, oauthTumblr := config.Socials["tumblr"]

		c.Set("oauth", oauthTwitter || oauthGitHub || oauthFacebook || oauthGoogle || oauthTumblr)
		c.Set("oauthTwitter", oauthTwitter)
		c.Set("oauthGitHub", oauthGitHub)
		c.Set("oauthFacebook", oauthFacebook)
		c.Set("oauthGoogle", oauthGoogle)
		c.Set("oauthTumblr", oauthTumblr)

		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}

func Login(c *gin.Context) {
	response := Response{} // todo sync.Pool
	defer response.Recover(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}
	// clean errors from client
	response.CleanErrors()

	// validate
	response.Username = strings.ToLower(response.Username)
	if len(response.Username) == 0 {
		response.ErrFor["username"] = "required"
	}
	if len(response.Password) == 0 {
		response.ErrFor["password"] = "required"
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()

	// abuseFilter
	collection := db.C(LOGINATTEMPTS)
	IpCountChan := make(chan int)
	IpUserCountChan := make(chan int)
	clientIP := c.ClientIP()
	go getCount(collection, IpCountChan, bson.M{
		"ip": clientIP,
	})
	go getCount(collection, IpUserCountChan, bson.M{
		"ip": clientIP,
		"user": response.Username,
	})
	IpCount := <- IpCountChan
	IpUserCount := <- IpUserCountChan
	if IpCount > config.LoginAttempts.ForIp || IpUserCount > config.LoginAttempts.ForIpAndUser {
		response.Errors = append(response.Errors, "You've reached the maximum number of login attempts. Please try again later.")
		response.Fail(c)
		return
	}

	// attemptLogin
	collection = db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"$or": []bson.M{
		bson.M{"username": response.Username},
		bson.M{"email": response.Email},
	}}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			response.Errors = append(response.Errors, "check username and password")
			response.Fail(c)
			return
		}
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(response.Password))
	if err != nil {
		attempt := LoginAttempt{}
		attempt.ID = bson.NewObjectId()
		attempt.IP = clientIP
		attempt.User = response.Username
		collection = db.C(LOGINATTEMPTS)
		err = collection.Insert(attempt)
		if err != nil {
			panic(err)
		}
		response.Errors = append(response.Errors, "check username and password")
		response.Fail(c)
		return
	}

	session := sessions.Default(c)
	session.Set("public", user.ID.Hex())
	if returnURL, ok := session.Get("returnURL").(string); ok {
		c.Redirect(http.StatusFound, returnURL)
	}
	session.Delete("returnURL")
	session.Save()

	response.Success = true
	c.JSON(http.StatusOK, response)
}
