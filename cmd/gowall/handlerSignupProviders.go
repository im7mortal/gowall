package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/markbates/goth"
	"gopkg.in/mgo.v2"
	"html/template"
)

func SignUpProvider(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("action", "/signup/")
	session.Save()
	startOAuth(c)
}

func signupProvider(c *gin.Context, userGoth goth.User) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err := collection.Find(bson.M{userGoth.Provider + ".id": userGoth.UserID}).One(&user)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		session := sessions.Default(c)
		session.Set("oauthMessage", "We found a user linked to your " + userGoth.Provider + " account")
		session.Save()
		c.Redirect(http.StatusFound, "/signup/")
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	userGothString, err := json.Marshal(userGoth)
	if err != nil {
		panic(err)
	}

	session := sessions.Default(c)
	session.Set("socialProfile", string(userGothString))
	session.Set("provider", userGoth.Provider)
	session.Save()

	c.Set("email", template.JS(userGoth.Email))
	c.HTML(http.StatusOK, "/signup/social/", c.Keys)
}

func SignUpSocial(c *gin.Context) {
	response := responseUser{}
	response.Init(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
	}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)

	// validate
	response.ValidateEmail(&response.Response)
	if response.HasErrors() {
		response.Fail()
		return
	}

	// check duplicate
	session := sessions.Default(c)

	socialProfile_, ok := session.Get("socialProfile").(string)
	if !ok || len(socialProfile_) == 0 {
		response.Errors = append(response.Errors, "something went wrong. Refresh please")
		response.Fail()
		return
	}
	socialProfile := goth.User{}
	err = json.Unmarshal([]byte(socialProfile_), &socialProfile)
	if err != nil {
		panic(err)
	}

	// duplicateEmailCheck
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	println(response.Email)
	err = collection.Find(bson.M{"email": response.Email}).One(&user)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.ErrFor["email"] = "email already registered"
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	// duplicateUsernameCheck
	var username string
	if len(socialProfile.Name) != 0 {
		username = socialProfile.Name
	} else if len(socialProfile.UserID) != 0 {
		username = socialProfile.UserID
	}
	reg, err := regexp.Compile(`/[^a-zA-Z0-9\-\_]/g`)
	if err != nil {
		panic(err)
	}
	usernameSrc := []byte(username)
	reg.ReplaceAll(usernameSrc, []byte(""))
	username = string(usernameSrc)
	if len(user.Username) != 0 {
		response.Fail()
	}
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err == nil {
		username += "-gowallUser"
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// createUser
	user.ID = bson.NewObjectId()
	user.IsActive = "yes"
	user.Username = user.ID.Hex()
	user.Email = strings.ToLower(response.Email)
	user.Search = []string{username, response.Email}
	user.updateProvider(socialProfile)
	err = collection.Insert(user)
	if err != nil {
		panic(err)
		return
	}

	// createAccount
	account := Account{}

	account.ID = bson.NewObjectId()

	user.Roles.Account = account.ID

	err = collection.UpdateId(user.ID, user)
	if err != nil {
		panic(err)
		return
	}

	if config.RequireAccountVerification {
		account.IsVerified = "no"
	} else {
		account.IsVerified = "yes"
	}
	account.Name.Full = username
	account.User.ID = user.ID
	account.User.Name = user.Username
	account.Search = []string{username}

	collection = db.C(ACCOUNTS)
	err = collection.Insert(account)
	if err != nil {
		panic(err)
		return
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