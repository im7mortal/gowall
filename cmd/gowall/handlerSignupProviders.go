package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"gopkg.in/mgo.v2"
)


func init()  {
	gothic.Store = store
}

func startOAuth(c *gin.Context) {
	// don't like that hack
	// gothic was written for another path
	// I just put provider query
	provider := c.Param("provider")
	c.Request.URL.RawQuery += "provider=" + provider
	_, err := goth.GetProvider(provider)
	if err != nil {
		// TODO HACK
		redir := "http://" + c.Request.Host + "/signup_/facebook/callback"
		goth.UseProviders(
			facebook.New(config.Socials["facebook"].Key, config.Socials["facebook"].Secret, redir),
		)
	}
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func CompleteUserAuth(c *gin.Context) {
	// gothic was written for another path
	// i just put provider query
	provider := c.Param("provider")
	c.Request.URL.RawQuery += "&provider=" + provider
	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	userGoth, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		render, _ := TemplateStorage["/signup/"]
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"facebook.id": userGoth.UserID}).One(&user)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		render, _ := TemplateStorage["/signup/"]
		c.Set("oauthMessage", "We found a user linked to your " + provider + " account")
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	userGothString, err := json.Marshal(userGoth)
	if err != nil {
		panic(err)
	}

	sessionCookie := sessions.Default(c)
	sessionCookie.Set("socialProfile", userGothString)
	sessionCookie.Set("provider", provider)
	sessionCookie.Save()

	c.Set("email", userGoth.Email)
	render, _ := TemplateStorage["/signup/social/"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func SignUpSocial(c *gin.Context) {
	response := Response{}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
	}

	// validate
	response.ValidateEmail()
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	// check duplicate
	session := sessions.Default(c)

	socialProfile_, ok := session.Get("socialProfile").(string)
	if !ok || len(socialProfile_) == 0 {
		render, _ := TemplateStorage["/signup/"]
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
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
	err = collection.Find(bson.M{"email": response.Email}).One(&user)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		render, _ := TemplateStorage["/signup/"]
		c.Set("oauthMessage", "email already registered") // TODO
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
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
		response.Fail(c)
	}
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err == nil {
		render, _ := TemplateStorage["/signup/"]
		c.Set("oauthMessage", "email already registered") // TODO
		render.Data = c.Keys
		c.Render(http.StatusOK, render)
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// createUser
	user.ID = bson.NewObjectId()
	user.IsActive = "yes"
	user.Username = user.ID.Hex()
	user.Email = strings.ToLower(response.Email)
	user.Search = []string{username, response.Email}
	user.Facebook = vendorOauth{}
	user.Facebook.ID = socialProfile.UserID
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

	response.Success = true
	c.JSON(http.StatusOK, response)
}