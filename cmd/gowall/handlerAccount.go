package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/contrib/sessions"
	"encoding/json"
	"net/url"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
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
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func AccountVerificationRender(c *gin.Context) {
	account, _ := getAccount(c)
	user, _ := getUser(c)
	if account.IsVerified == "yes" {
		c.HTML(http.StatusOK, user.DefaultReturnUrl(), c.Keys)
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

		session, err := mgo.Dial("mongodb://localhost:27017")
		defer session.Close()
		if err != nil {
			println(err.Error())
		}
		d := session.DB("test")
		collection := d.C(ACCOUNTS)
		account.VerificationToken = string(hash)
		collection.UpdateId(account.ID, account)// todo how to update only part?
		verifyURL := "http" +"://"+ "localhost:3000" +"/account/verification/" + string(VerifyURL) + "/"
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

	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func AccountSettingsRender(c *gin.Context) {
	sess := sessions.Default(c)

	public := sess.Get("public")
	session, err := mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	collection := d.C(USERS)
	us := User{} // todo pool
	err = collection.FindId(bson.ObjectIdHex(public.(string))).One(&us)
	if err != nil {
		println(err.Error())
	}
	if len(us.Username) != 0 {
		User, _ := json.Marshal(gin.H{
			"_id": us.ID.Hex(),
			"user": us.Username,
			"email": us.Email,
		})
		c.Set("User", url.QueryEscape(string(User)))
	}
	collection = d.C(ACCOUNTS)
	ac := Account{} // todo pool
	err = collection.FindId(bson.ObjectIdHex(us.Roles.Account.Hex())).One(&ac)
	if err != nil {
		println(err.Error())
	}
	if len(ac.ID) != 0 {
		Account, _ := json.Marshal(gin.H{
			"_id": ac.ID.Hex(),
			"name": gin.H{
				"first": ac.Name.First,
				"middle": ac.Name.Middle,
				"last": ac.Name.Last,
			},
			"company": ac.Company,
			"phone": ac.Phone,
			"zip": ac.Zip,
		})
		c.Set("Account", url.QueryEscape(string(Account)))
	}
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}