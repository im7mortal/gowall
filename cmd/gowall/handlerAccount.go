package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func renderAccount(c *gin.Context) {
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func renderAccountVerification(c *gin.Context) {
	account := getAccount(c)
	user := getUser(c)
	if account.IsVerified == "yes" {
		c.Redirect(http.StatusFound, user.defaultReturnUrl())
		return
	}
	if len(account.VerificationToken) > 0 {

	} else {
		VerifyURL := generateToken(21)
		hash, err := bcrypt.GenerateFromPassword(VerifyURL, bcrypt.DefaultCost)
		if err != nil {
			EXCEPTION(err)
		}
		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ACCOUNTS)
		err = collection.UpdateId(account.ID, bson.M{
			"$set": bson.M{
				"verificationToken": string(hash),
			},
		})
		if err != nil {
			EXCEPTION(err)
		}
		verifyURL := "http" + "://" + c.Request.Host + "/account/verification/" + string(VerifyURL) + "/"
		c.Set("VerifyURL", verifyURL)

		mailConf := MailConfig{}
		mailConf.Data = c.Keys
		mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
		mailConf.To = user.Email
		mailConf.Subject = "Your " + config.ProjectName + " Account"
		mailConf.ReplyTo = user.Email
		mailConf.HtmlPath = "views/account/verification/email-html.html"

		if err := mailConf.SendMail(); err != nil {
			EXCEPTION(err)
		}
	}

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func verify(c *gin.Context) {
	account := getAccount(c)
	user := getUser(c)
	err := bcrypt.CompareHashAndPassword([]byte(account.VerificationToken), []byte(c.Param("token")))
	if err == nil {
		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ACCOUNTS)
		collection.UpdateId(account.ID, bson.M{
			"$set": bson.M{
				"verificationToken": "",
				"isVerified":        "yes",
			},
		})
	}
	c.Redirect(http.StatusFound, user.defaultReturnUrl())
}

func resendVerification(c *gin.Context) {
	account := getAccount(c)
	user := getUser(c)
	if account.IsVerified == "yes" {
		c.HTML(http.StatusOK, user.defaultReturnUrl(), c.Keys)
		return
	}
	response := Response{}
	response.Init(c)

	var body struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&body)

	validateEmail(&body.Email, &response)

	if response.HasErrors() {
		response.Fail()
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	collection.Find(
		bson.M{
			"email": body.Email,
			"_id": bson.M{
				"$ne": user.ID,
			},
		}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "That email already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}
	collection.UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"email": body.Email,
		},
	})

	collection = db.C(ACCOUNTS)
	VerifyURL := generateToken(21)
	hash, err := bcrypt.GenerateFromPassword(VerifyURL, bcrypt.DefaultCost)
	if err != nil {
		EXCEPTION(err)
	}
	collection.UpdateId(account.ID, bson.M{
		"verificationToken": string(hash),
	})
	verifyURL := "http" + "://" + c.Request.Host + "/account/verification/" + string(VerifyURL) + "/"
	c.Set("VerifyURL", verifyURL)
	mailConf := MailConfig{}
	mailConf.Data = c.Keys
	mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
	mailConf.To = user.Email
	mailConf.Subject = "Your " + config.ProjectName + " Account"
	mailConf.ReplyTo = user.Email
	mailConf.HtmlPath = "views/account/verification/email-html.html"

	if err := mailConf.SendMail(); err != nil {
		EXCEPTION(err)
	}
	response.Finish()
}
