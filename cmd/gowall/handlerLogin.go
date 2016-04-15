package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/contrib/sessions"
	"regexp"
	"time"
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

		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}

func ForgotRender(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
	} else {
		render, _ := TemplateStorage[c.Request.URL.Path]
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}

func SendReset(c *gin.Context) {
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
	if response.HasErrors() {
		response.Fail(c)
		return
	}


	token := generateToken(21)
	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		println(err.Error())
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	us := User{} // todo pool
	err = collection.Find(bson.M{"email": email}).One(&us)
	if err != nil {
		println(err.Error())
	}
	if len(us.Username) == 0 {
		response.ErrFor["email"] = "email doesn't exist"
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	us.ResetPasswordToken = string(hash)
	us.ResetPasswordExpires = time.Now().Add(24 * time.Hour)
	collection.UpdateId(us.ID, us)

	resetURL := "http" +"://"+ "localhost:3000" +"/login/reset/" + email + "/" + string(token) + "/"
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
		//todo it's not serious
	}
	response.Success = true
	c.JSON(http.StatusOK, response)
}






func ResetRender(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		defaultReturnUrl, _ := c.Get("DefaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
	} else {
		render, _ := TemplateStorage["/login/reset/"] // can't handle /login/reset/:email:token
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}

func ResetPassword (c *gin.Context) {

	var body struct {
		Confirm   string  `json:"confirm"`
		Password string  `json:"password"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)


	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)


	password := strings.ToLower(body.Password)
	if len(password) == 0 {
		response.ErrFor["password"] = "required"
	}
	confirm := strings.ToLower(body.Confirm)
	if len(confirm) == 0 {
		response.ErrFor["confirm"] = "required"
	}
	if confirm != password {
		response.Errors = append(response.Errors,"Passwords do not match.")
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	us := User{}
	err = collection.Find(bson.M{"email": c.Param("email"), "resetPasswordExpires": bson.M{"$gt": time.Now()}}).One(&us)

	if err != nil {
		println(err.Error())
		response.Fail(c)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(us.ResetPasswordToken), []byte(c.Param("token")))

	if err == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
			response.Fail(c)
			return
		}
		us.Password = string(hashedPassword)
		collection.UpdateId(us.ID, us)
	}
	response.Success = true
	c.JSON(http.StatusOK, response)
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


	// TODO  abuseFilter!!!!!!!!!!!!!
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(LOGINATTEMPTS)
	at := LoginAttempt{} // todo pool
	at.ID = bson.NewObjectId()
	at.IP = c.ClientIP()
	at.User = username
	err = collection.Insert(at)
	if err != nil {
		println(err.Error())
	}

	collection = db.C(USERS)
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
		if returnURL_, ok_ := sess.Get("returnURL").(string); ok_ {
			returnURL = returnURL_
		}
		sess.Delete("returnURL")
		sess.Save()
	}

	if len(returnURL) > 0 {
		c.Redirect(http.StatusFound, returnURL)
	} else {
		c.JSON(http.StatusOK, response)
	}
}