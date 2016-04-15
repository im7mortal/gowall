package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/contrib/sessions"
	"encoding/json"
	"net/url"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

func generateToken1(n int) []byte {
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

func AdminRender(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]

	// TODO
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	count, _ := collection.Count()
	c.Set("CountAccount", count)
	collection = db.C(USERS)
	count, _ = collection.Count()
	c.Set("CountUser", count)
	collection = db.C(ADMINS)
	count, _ = collection.Count()
	c.Set("CountAdmin", count)
	collection = db.C(ADMINGROUPS)
	count, _ = collection.Count()
	c.Set("CountAdminGroup", count)
	collection = db.C(CATEGORIES)
	count, _ = collection.Count()
	c.Set("CountCategory", count)
	collection = db.C(STATUS)
	count, _ = collection.Count()
	c.Set("CountStatus", count)

	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func AccountVerificationRender1(c *gin.Context) {
	account, _ := getAccount(c)
	user, _ := getUser(c)
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

func Verify1 (c *gin.Context) {

	account, _ := getAccount(c)
	user, _ := getUser(c)
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

func ResendVerification1 (c *gin.Context) {
	account, _ := getAccount(c)
	user, _ := getUser(c)
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

	response.Success = true
	c.JSON(http.StatusOK, response)
}

func AccountSettingsRender1(c *gin.Context) {
	sess := sessions.Default(c)

	public := sess.Get("public")
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	us := User{} // todo pool
	err := collection.FindId(bson.ObjectIdHex(public.(string))).One(&us)
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
	collection = db.C(ACCOUNTS)
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