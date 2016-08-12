package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"gopkg.in/mgo.v2"
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
	response := newResponse(c)

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}

	// validate
	validateUsername(&body.Username, response)
	validateEmail(&body.Email, response)
	validatePassword(&body.Password, response)

	if response.HasErrors() {
		response.Fail()
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.Find(
		bson.M{"$or": []bson.M{
			bson.M{"username": body.Username},
			bson.M{"email": body.Email},
		},
		}).One(&user)

	if err != nil && err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// duplicateUsernameCheck
	// duplicateEmailCheck
	if len(user.Username) != 0 {
		if user.Username == body.Username {
			response.ErrFor["username"] = "username already taken."
		}
		if user.Email == body.Email {
			response.ErrFor["email"] = "email already registered."
		}
	}
	if response.HasErrors() {
		response.Fail()
		return
	}

	// createUser
	user.setPassword(body.Password)

	user.ID = bson.NewObjectId()
	user.IsActive = "yes"
	user.Username = body.Username
	user.Email = strings.ToLower(body.Email)
	user.Search = []string{body.Username, body.Email}

	// createAccount
	account := Account{}
	account.ID = bson.NewObjectId()

	// insertUser
	user.Roles.Account = account.ID
	err = collection.Insert(user)
	if err != nil {
		EXCEPTION(err)
	}

	// insertAccount
	if config.RequireAccountVerification {
		account.IsVerified = "no"
	} else {
		account.IsVerified = "yes"
	}
	account.Name.Full = body.Username
	account.User.ID = user.ID
	account.User.Name = user.Username
	account.Search = []string{body.Username}

	collection = db.C(ACCOUNTS)
	err = collection.Insert(account)
	if err != nil {
		EXCEPTION(err)
	}

	// sendWelcomeEmail
	go func() {
		c.Set("Username", body.Username)
		c.Set("Email", body.Email)
		c.Set("LoginURL", "http://" + c.Request.Host + "/login/")

		mailConf := MailConfig{}
		mailConf.Data = c.Keys
		mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
		mailConf.To = config.SystemEmail
		mailConf.Subject = "Your " + config.ProjectName + " Account"
		mailConf.ReplyTo = body.Email
		mailConf.HtmlPath = "views/signup/email-html.html"

		if err := mailConf.SendMail(); err != nil {
			println("Error Sending Welcome Email: " + err.Error())
		}
	}()

	// logUserIn
	user.login(c)

	response.Finish()
}
