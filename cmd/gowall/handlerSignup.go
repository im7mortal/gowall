package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

func renderSignup(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		defaultReturnUrl, exist := c.Get("DefaultReturnUrl")
		var url string
		if url, ok = defaultReturnUrl.(string); !exist || !ok {
			// if not exist or not string
			url = "/"
		}
		c.Redirect(http.StatusFound, url)
	} else {
		injectSocials(c)
		c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
	}
}

func signup(c *gin.Context) {
	response := responseUser{}
	response.Init(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		EXCEPTION(err)
	}
	// clean errors from client
	response.CleanErrors()

	// validate
	validateUsername(&response.User.Username, &response.Response)
	validateEmail(&response.User.Email, &response.Response)
	validatePassword(&response.User.Password, &response.Response)

	if response.HasErrors() {
		response.Fail()
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"$or": []bson.M{bson.M{"username": response.Username}, bson.M{"email": response.Email}}}).One(&user)
	if err != nil {
		EXCEPTION(err)
	}

	// duplicateUsernameCheck
	// duplicateEmailCheck
	if len(user.Username) != 0 {
		if user.Username == response.Username {
			response.ErrFor["username"] = "username already taken"
		}
		if user.Email == response.Email {
			response.ErrFor["email"] = "email already registered"
		}
	}
	if response.HasErrors() {
		response.Fail()
		return
	}

	// createUser
	user.setPassword(response.Password)

	user.ID = bson.NewObjectId()
	user.IsActive = "yes"
	user.Username = response.Username
	user.Email = strings.ToLower(response.Email)
	user.Search = []string{response.Username, response.Email}

	err = collection.Insert(user)
	if err != nil {
		EXCEPTION(err)
	}

	// createAccount
	account := Account{}

	account.ID = bson.NewObjectId()

	//update user with account
	user.Roles.Account = account.ID
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		EXCEPTION(err)
	}

	if config.RequireAccountVerification {
		account.IsVerified = "no"
	} else {
		account.IsVerified = "yes"
	}
	account.Name.Full = response.Username
	account.User.ID = user.ID
	account.User.Name = user.Username
	account.Search = []string{response.Username}

	collection = db.C(ACCOUNTS)
	err = collection.Insert(account)
	if err != nil {
		EXCEPTION(err)
	}

	// sendWelcomeEmail
	go func() {
		c.Set("Username", response.Username)
		c.Set("Email", response.Email)
		c.Set("LoginURL", "http://" + c.Request.Host + "/login/")

		mailConf := MailConfig{}
		mailConf.Data = c.Keys
		mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
		mailConf.To = config.SystemEmail
		mailConf.Subject = "Your " + config.ProjectName + " Account"
		mailConf.ReplyTo = response.Email
		mailConf.HtmlPath = "views/signup/email-html.html"

		if err := mailConf.SendMail(); err != nil {
			println("Error Sending Welcome Email: " + err.Error())
		}
	}()

	// logUserIn
	user.login(c)

	response.Finish()
}
