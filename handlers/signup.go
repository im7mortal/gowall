package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"github.com/im7mortal/gowall/config"
	"github.com/im7mortal/gowall/schemas"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/contrib/sessions"
)

func SignupRender(c *gin.Context) {
	_, isAuthenticated := c.Get("isAuthenticated") // non standard way. If exist it isAuthenticated
	if isAuthenticated {
		defaultReturnUrl, _ := c.Get("defaultReturnUrl")
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

func Signup(c *gin.Context) {
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
	} else {
		r, err := regexp.MatchString(`^[a-zA-Z0-9\-\_]+$`, username)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
	email := strings.ToLower(body.Email)
	if len(email) == 0 {
		response.ErrFor["email"] = "required"
	} else {
		r, err := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\+]+@[a-zA-Z0-9\-\_\.]+\.[a-zA-Z0-9\-\_]+$`, email)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["email"] = `invalid email format`
		}
	}
	password := body.Password
	//destruct body??? todo
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
	d := session.DB("test")
	collection := d.C("User")
	//collection := d.C("restaurants")
	collection.Create(&mgo.CollectionInfo{})
	us := schemas.User{} // todo pool
	err = collection.Find(bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}).One(&us)
	if err != nil {
		println(err.Error())
	}
	if len(us.Username) != 0 {
		if us.Username == username {
			response.ErrFor["username"] = "username already taken"
		}
		if us.Email == email {
			response.ErrFor["email"] = "email already registered"
		}
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	us.ID = bson.NewObjectId()
	us.IsActive = "yes"
	us.Username = username
	us.Email = strings.ToLower(email)
	us.Password = string(hashedPassword)
	us.Search = []string{username, email}

	err = collection.Insert(us)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	collection = d.C("Account")
	//collection := d.C("restaurants")
	collection.Create(&mgo.CollectionInfo{})

	ac := schemas.Account{}

	ac.ID = bson.NewObjectId()
	if config.RequireAccountVerification {
		ac.IsVerified = "no"
	} else {
		ac.IsVerified = "yes"
	}
	ac.Name.Full = username
	ac.User.ID = us.ID
	ac.User.Name = us.Username
	ac.Search = []string{username}
	err = collection.Insert(ac)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	//todo  sendWelcomeEmail
	//todo  sendWelcomeEmail ***************************************************
	sess := sessions.Default(c)
	sess.Set("public", us.ID.Hex())
	sess.Save()

	response.Success = false
	c.JSON(http.StatusOK, response)
}