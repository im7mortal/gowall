package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/contrib/sessions"
)

func LoginRender(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		var redirectURL string
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		redirectURL = defaultReturnUrl.(string)
		sess := sessions.Default(c)
		returnURL := sess.Get("returnURL")
		if returnURL != nil {
			redirectURL = returnURL.(string)
			sess.Delete("returnURL")
			sess.Save()
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
		c.Set("oauthMessage", "")

		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}

func ForgotRender(c *gin.Context) {
	_, isAuthenticated := c.Get("isAuthenticated") // non standard way. If exist it isAuthenticated
	if isAuthenticated {
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
	} else {
		render, _ := TemplateStorage[c.Request.URL.Path]
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}

func Login(c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)

	var body struct {
		Username    string  `json:"username"`
		Email   string  `json:"email"`
		Password string  `json:"password"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	username := strings.ToLower(body.Username)
	if len(username) == 0 {
		response.ErrFor["username"] = "required"
	}
	password := body.Password
	if len(password) == 0 {
		response.ErrFor["password"] = "required"
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}

	// TODO  abuseFilter!!!!!!!!!!!!!
	session, err = mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	collection := d.C(LOGINATTEMPTS)
	collection.Create(&mgo.CollectionInfo{})
	at := LoginAttempt{} // todo pool
	at.ID = bson.NewObjectId()
	at.IP = c.ClientIP()
	at.User = username
	err = collection.Insert(at)
	if err != nil {
		println(err.Error())
	}

	collection = d.C(USERS)
	us := User{}
	err = collection.Find(bson.M{"username": username}).One(&us)
	if err != nil {
		println(err.Error())
	}
	var returnURL string
	if len(us.Username) > 0 {
		err := bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(password))
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			response.Fail(c)
			return
		}
		sess := sessions.Default(c)
		sess.Set("public", us.ID.Hex())
		response.Success = true
		returnURL = sess.Get("returnURL").(string)
		sess.Delete("returnURL")
		sess.Save()
	}

	if len(returnURL) > 0 {
		c.Redirect(http.StatusFound, returnURL)
	} else {
		c.JSON(http.StatusOK, response)
	}
}