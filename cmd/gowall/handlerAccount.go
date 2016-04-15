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
	"regexp"
"strings"
	"html/template"
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

		session, err := mgo.Dial(config.MongoDB)
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

func Verify (c *gin.Context) {

	account, _ := getAccount(c)
	user, _ := getUser(c)
	err := bcrypt.CompareHashAndPassword([]byte(account.VerificationToken), []byte(c.Param("token")))
	if err == nil {
		session, err := mgo.Dial(config.MongoDB)
		defer session.Close()
		if err != nil {
			println(err.Error())
		}
		d := session.DB("test")
		collection := d.C(ACCOUNTS)
		account.VerificationToken = ""
		account.IsVerified = "yes"
		collection.UpdateId(account.ID, account)
	}
	c.Redirect(http.StatusFound, user.DefaultReturnUrl())
}

func ResendVerification (c *gin.Context) {
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
	session, err := mgo.Dial(config.MongoDB)
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	collection := d.C(USERS)
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

	collection = d.C(ACCOUNTS)
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

func AccountSettingsRender(c *gin.Context) {
	sess := sessions.Default(c)

	public := sess.Get("public")
	session, err := mgo.Dial(config.MongoDB)
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
		c.Set("User", template.JS(url.QueryEscape(string(User))))
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
		c.Set("Account", template.JS(url.QueryEscape(string(Account))))
	}
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func SetSettings (c *gin.Context) {
	account, _ := getAccount(c)
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)
	var body struct {
		First   string  `json:"first"`
		Middle  string  `json:"middle"`
		Last    string  `json:"last"`
		Company string  `json:"company"`
		Phone   string  `json:"phone"`
		Zip     string  `json:"zip"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	if len(body.First) == 0 {
		response.ErrFor["first"] = "required"
	}
	if len(body.Last) == 0 {
		response.ErrFor["last"] = "required"
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}
	session, err := mgo.Dial(config.MongoDB)
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")

	account.Name.Full = body.First + " " + body.Last
	account.Name.First = body.First
	account.Name.Middle = body.Middle
	account.Name.Last = body.Last
	account.Company = body.Company
	account.Phone = body.Phone
	account.Zip = body.Zip
	account.Search = account.Search[:0]
	account.Search = append(account.Search,
		body.First,
		body.Middle,
		body.Last,
		body.Company,
		body.Phone,
		body.Zip,
	)

	collection := d.C(ACCOUNTS)
	err = collection.UpdateId(account.ID, account)
	if err != nil {
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

func ChangePassword (c *gin.Context) {
	user, _ := getUser(c)
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)
	var body struct {
		Confirm   string  `json:"confirm"`
		Password string  `json:"newPassword"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	if len(body.Password) == 0 {
		response.ErrFor["newPassword"] = "required"
	}
	if len(body.Confirm) == 0 {
		response.ErrFor["confirm"] = "required"
	} else if body.Password != body.Confirm {
		response.Errors = append(response.Errors, "Passwords do not match.")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}
	session, err := mgo.Dial(config.MongoDB)
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Errors = append(response.Errors, err.Error()) // TODO don't like that this error goes to client
		response.Fail(c)
		return
	}

	user.Password = string(hashedPassword)
	collection := d.C(USERS)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

func ChangeIdentity (c *gin.Context) {
	user, _ := getUser(c)
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)
	var body struct {
		Username    string  `json:"username"`
		Email   string  `json:"email"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)


	username := strings.ToLower(body.Username)
	if len(username) == 0 {
		response.ErrFor["username"] = "required"
	} else {
		r, err := regexp.MatchString(`^[a-zA-Z0-9\-\_]+$`, username)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
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
	session, err := mgo.Dial(config.MongoDB)
	defer session.Close()
	if err != nil {
		println(err.Error())
	}

	d := session.DB("test")
	collection := d.C(USERS)

	{
		us := User{} // todo pool
		err = collection.Find(bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}).One(&us)
		if err != nil {
			response.Errors = append(response.Errors, "username or email already exist")
			response.Fail(c)
			return
		}
	}

	user.Username = body.Username
	user.Email = body.Email

	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}
	// TODO  patch admin and account
	response.Success = true
	c.JSON(http.StatusOK, response)
}
