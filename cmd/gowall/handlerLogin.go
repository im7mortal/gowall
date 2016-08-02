package main

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
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
		injectSocials(c)
		c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
	}
}

func Login(c *gin.Context) {
	response := responseUser{}
	response.Init(c)

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}
	// clean errors from client
	response.CleanErrors()

	// validate
	response.Username = strings.ToLower(response.Username)
	if len(body.Username) == 0 {
		response.ErrFor["username"] = "required"
	}
	if len(body.Password) == 0 {
		response.ErrFor["password"] = "required"
	}
	if response.HasErrors() {
		response.Fail()
		return
	}
	response.Username = body.Username
	response.Password = body.Password
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
		"ip":   clientIP,
		"user": response.Username,
	})
	IpCount := <-IpCountChan
	IpUserCount := <-IpUserCountChan
	if IpCount > config.LoginAttempts.ForIp || IpUserCount > config.LoginAttempts.ForIpAndUser {
		response.Errors = append(response.Errors, "You've reached the maximum number of login attempts. Please try again later.")
		response.Fail()
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
			response.Fail()
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
		response.Fail()
		return
	}

	session := sessions.Default(c)
	session.Set("public", user.ID.Hex())
	if returnURL, ok := session.Get("returnURL").(string); ok {
		c.Redirect(http.StatusFound, returnURL)
	}
	session.Delete("returnURL")
	session.Save()

	response.Finish()
}
