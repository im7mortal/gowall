package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/contrib/sessions"
)

func signupRender(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		defaultReturnUrl, exist := c.Get("DefaultReturnUrl")
		var url string
		if url, ok = defaultReturnUrl.(string); !exist || !ok { // if not exist or not string
			url = "/"
		}
		c.Redirect(http.StatusFound, url)
	} else {
		injectSocials(c)
		c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
	}
}

func Signup(c *gin.Context) {
	response := responseUser{}
	response.Init(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
		return
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
		println(err.Error())
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(response.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
		return
	}

	user.ID = bson.NewObjectId()
	user.IsActive = "yes"
	user.Username = response.Username
	user.Email = strings.ToLower(response.Email)
	user.Password = string(hashedPassword)
	user.Search = []string{response.Username, response.Email}

	err = collection.Insert(user)
	if err != nil {
		panic(err)
		return
	}

	// createAccount
	account := Account{}

	account.ID = bson.NewObjectId()

	//update user with account
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
	account.Name.Full = response.Username
	account.User.ID  = user.ID
	account.User.Name = user.Username
	account.Search = []string{response.Username}

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

func (user *User)login(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Set("public", user.ID.Hex())
	sess.Save()
}
