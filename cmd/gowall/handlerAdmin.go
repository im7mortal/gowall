package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
)

func generateToken1(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		println(err.Error())
		return b
	}
	token := make([]byte, n*2)
	hex.Encode(token, b)
	return token
}

func renderAdministrator(c *gin.Context) {

	// todo has to be sync pkg
	db := getMongoDBInstance()
	defer db.Session.Close()

	CountAccountChan := make(chan int)
	CountUserChan := make(chan int)
	CountAdminChan := make(chan int)
	CountAdminGroupChan := make(chan int)
	CountCategoryChan := make(chan int)
	CountStatusChan := make(chan int)

	go getCount(db.C(ACCOUNTS), CountAccountChan, bson.M{})
	go getCount(db.C(USERS), CountUserChan, bson.M{})
	go getCount(db.C(ADMINS), CountAdminChan, bson.M{})
	go getCount(db.C(ADMINGROUPS), CountAdminGroupChan, bson.M{})
	go getCount(db.C(CATEGORIES), CountCategoryChan, bson.M{})
	go getCount(db.C(STATUSES), CountStatusChan, bson.M{})

	c.Set("CountAccount", <-CountAccountChan)
	c.Set("CountUser", <-CountUserChan)
	c.Set("CountAdmin", <-CountAdminChan)
	c.Set("CountAdminGroup", <-CountAdminGroupChan)
	c.Set("CountCategory", <-CountCategoryChan)
	c.Set("CountStatus", <-CountStatusChan)

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func AccountVerificationRender1(c *gin.Context) {
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
			println(err.Error())
			return
		}

		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ACCOUNTS)
		account.VerificationToken = string(hash)
		collection.UpdateId(account.ID, account) // todo how to update only part?
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
			//todo it's not serious
		}

	}

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func Verify1(c *gin.Context) {

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
	c.Redirect(http.StatusFound, user.defaultReturnUrl())
}

func ResendVerification1(c *gin.Context) {
	account := getAccount(c)
	user := getUser(c)
	if account.IsVerified == "yes" {
		c.HTML(http.StatusOK, user.defaultReturnUrl(), c.Keys)
		return
	}
	response := Response{}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	validateEmail(&body.Email, &response)
	if response.HasErrors() {
		response.Fail()
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	{
		user_ := User{}
		collection.Find(bson.M{"email": body.Email, "_id": bson.M{"$ne": user.ID}}).One(&user_)
		if len(user_.Username) > 0 {
			response.ErrFor["email"] = `email already taken`
		}
	}
	if response.HasErrors() {
		response.Fail()
		return
	}
	user.Email = body.Email
	collection.UpdateId(user.ID, user)

	collection = db.C(ACCOUNTS)
	VerifyURL := generateToken(21)
	hash, err := bcrypt.GenerateFromPassword(VerifyURL, bcrypt.DefaultCost)
	if err != nil {
		println(err.Error())
		return
	}
	account.VerificationToken = string(hash)
	collection.UpdateId(account.ID, account) // todo how to update only part?
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
		//todo it's not serious
	}

	response.Finish()
}

func AccountSettingsRender1(c *gin.Context) {
	sess := sessions.Default(c)

	public := sess.Get("public")
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	us := User{}
	err := collection.FindId(bson.ObjectIdHex(public.(string))).One(&us)
	if err != nil {
		println(err.Error())
	}
	if len(us.Username) != 0 {
		User, _ := json.Marshal(gin.H{
			"_id":   us.ID.Hex(),
			"user":  us.Username,
			"email": us.Email,
		})
		c.Set("User", url.QueryEscape(string(User)))
	}
	collection = db.C(ACCOUNTS)
	ac := Account{}
	err = collection.FindId(us.Roles.Account).One(&ac)
	if err != nil {
		println(err.Error())
	}
	if len(ac.ID) != 0 {
		Account, _ := json.Marshal(gin.H{
			"_id": ac.ID.Hex(),
			"name": gin.H{
				"first":  ac.Name.First,
				"middle": ac.Name.Middle,
				"last":   ac.Name.Last,
			},
			"company": ac.Company,
			"phone":   ac.Phone,
			"zip":     ac.Zip,
		})
		c.Set("Account", url.QueryEscape(string(Account)))
	}
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}
