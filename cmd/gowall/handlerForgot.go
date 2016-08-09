package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
	"gopkg.in/mgo.v2"
)

func renderForgot(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
	} else {
		c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
	}
}

func sendReset(c *gin.Context) {
	response := Response{}
	response.Init(c)

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}

	validateEmail(&body.Email, &response)
	if response.HasErrors() {
		response.Fail()
		return
	}

	token := generateToken(21)
	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		EXCEPTION(err)
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	us := User{}
	err = collection.Find(bson.M{"email": body.Email}).One(&us)
	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.ErrFor["email"] = "email doesn't exist."
	}
	if response.HasErrors() {
		response.Fail()
		return
	}

	us.ResetPasswordToken = string(hash)
	us.ResetPasswordExpires = time.Now().Add(24 * time.Hour)
	collection.UpdateId(us.ID, us)

	resetURL := "http" + "://" + c.Request.Host + "/login/reset/" + body.Email + "/" + string(token) + "/"
	c.Set("ResetURL", resetURL)
	c.Set("Username", us.Username)

	mailConf := MailConfig{}
	mailConf.Data = c.Keys
	mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
	mailConf.To = config.SystemEmail
	mailConf.Subject = "Your " + config.ProjectName + " Account"
	mailConf.ReplyTo = body.Email
	mailConf.HtmlPath = "views/login/forgot/email-html.html"

	if err := mailConf.SendMail(); err != nil {
		EXCEPTION(err)
	}
	response.Finish()
}

func renderReset(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
	} else {
		c.HTML(http.StatusOK, "/login/reset/", c.Keys) // can't handle /login/reset/:email:token
	}
}

func resetPassword(c *gin.Context) {

	var body struct {
		Confirm  string `json:"confirm"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	response := Response{}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)

	password := strings.ToLower(body.Password)
	if len(password) == 0 {
		response.ErrFor["password"] = "required"
	}
	confirm := strings.ToLower(body.Confirm)
	if len(confirm) == 0 {
		response.ErrFor["confirm"] = "required"
	}
	if confirm != password {
		response.Errors = append(response.Errors, "Passwords do not match.")
	}
	if response.HasErrors() {
		response.Fail()
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"email": c.Param("email"), "resetPasswordExpires": bson.M{"$gt": time.Now()}}).One(&user)

	if err != nil {
		response.Fail()
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.ResetPasswordToken), []byte(c.Param("token")))

	if err == nil {

		user.setPassword(password)
		collection.UpdateId(user.ID, user)
	}
	response.Finish()
}
