package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/im7mortal/gowall/config"
	"encoding/json"
	"strings"
	"gopkg.in/mgo.v2"
	"github.com/im7mortal/gowall/schemas"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/contrib/sessions"
)

func LoginRender(c *gin.Context) {
	_, isAuthenticated := c.Get("isAuthenticated") // non standard way. If exist it isAuthenticated
	if isAuthenticated {
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
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
	collection := d.C("LoginAttempt")
	collection.Create(&mgo.CollectionInfo{})
	at := schemas.LoginAttempt{} // todo pool
	at.ID = bson.NewObjectId()
	at.IP = c.ClientIP()
	at.User = username
	err = collection.Insert(at)
	if err != nil {
		println(err.Error())
	}

	collection = d.C("User")
	us := schemas.User{}
	err = collection.Find(bson.M{"username": username}).One(&us)
	if err != nil {
		println(err.Error())
	}
	if len(us.Username) > 0 {
		err := bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(password))
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			response.Fail(c)
			return
		}
		sess := sessions.Default(c)
		sess.Set("public", us.ID.Hex())
		sess.Save()
		response.Success = true
	}
	c.JSON(http.StatusOK, response)
}