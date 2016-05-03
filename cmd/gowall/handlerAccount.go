package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

func generateToken(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		println(err.Error())
		return b
	}
	token := make([]byte, n * 2)
	hex.Encode(token, b)
	return token
}

func AccountRender(c *gin.Context) {
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func AccountVerificationRender(c *gin.Context) {
	account := getAccount(c)
	user := getUser(c)
	if account.IsVerified == "yes" {
		c.Redirect(http.StatusFound, user.DefaultReturnUrl())
		return
	}
	if len(account.VerificationToken) > 0 {

	} else {
		VerifyURL := generateToken(21)
		hash, err := bcrypt.GenerateFromPassword(VerifyURL, bcrypt.DefaultCost)
		if err != nil {
			println(err.Error())
			return
		}

		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ACCOUNTS)
		account.VerificationToken = string(hash)
		collection.UpdateId(account.ID, account)// todo how to update only part?
		verifyURL := "http" +"://"+ c.Request.Host +"/account/verification/" + string(VerifyURL) + "/"
		c.Set("VerifyURL", verifyURL)

		mailConf := MailConfig{}
		mailConf.Data = c.Keys
		mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
		mailConf.To = user.Email
		mailConf.Subject = "Your " + config.ProjectName + " Account"
		mailConf.ReplyTo = user.Email
		mailConf.HtmlPath = "views/account/verification/email-html.html"

		if err := mailConf.SendMail(); err != nil {
			//todo it's not serious
		}



	}

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func Verify (c *gin.Context) {

	account := getAccount(c)
	user := getUser(c)
	err := bcrypt.CompareHashAndPassword([]byte(account.VerificationToken), []byte(c.Param("token")))
	if err == nil {
		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ACCOUNTS)
		account.VerificationToken = ""
		account.IsVerified = "yes"
		collection.UpdateId(account.ID, account)
	}
	c.Redirect(http.StatusFound, user.DefaultReturnUrl())
}

func ResendVerification (c *gin.Context) {
	account := getAccount(c)
	user := getUser(c)
	if account.IsVerified == "yes" {
		c.HTML(http.StatusOK, user.DefaultReturnUrl(), c.Keys)
		return
	}
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
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	{
		user_ := User{}
		collection.Find(bson.M{"email": email, "_id": bson.M{"$ne": user.ID}}).One(&user_)
		if len(user_.Username) > 0 {
			response.ErrFor["email"] = `email already taken`
		}
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}
	user.Email = email
	collection.UpdateId(user.ID, user)

	collection = db.C(ACCOUNTS)
	VerifyURL := generateToken(21)
	hash, err := bcrypt.GenerateFromPassword(VerifyURL, bcrypt.DefaultCost)
	if err != nil {
		println(err.Error())
		return
	}
	account.VerificationToken = string(hash)
	collection.UpdateId(account.ID, account)// todo how to update only part?
	verifyURL := "http" +"://"+ c.Request.Host +"/account/verification/" + string(VerifyURL) + "/"
	c.Set("VerifyURL", verifyURL)
	mailConf := MailConfig{}
	mailConf.Data = c.Keys
	mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
	mailConf.To = user.Email
	mailConf.Subject = "Your " + config.ProjectName + " Account"
	mailConf.ReplyTo = user.Email
	mailConf.HtmlPath = "views/account/verification/email-html.html"

	if err := mailConf.SendMail(); err != nil {
		//todo it's serious
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}
