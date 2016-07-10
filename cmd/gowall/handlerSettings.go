package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gin-gonic/contrib/sessions"
	"encoding/json"
	"html/template"
	"github.com/markbates/goth"
	"gopkg.in/mgo.v2"
	"regexp"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"net/url"
)

func AccountSettingsRender(c *gin.Context) {
	sess := sessions.Default(c)

	user := getUser(c)
	injectSocials(c)
	doUserHasSocials(c, user)

	public := sess.Get("public")
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user = &User{} // todo pool
	err := collection.FindId(bson.ObjectIdHex(public.(string))).One(user)
	if err != nil {
		println(err.Error())
	}
	if len(user.Username) != 0 {
		User, _ := json.Marshal(gin.H{
			"_id": user.ID.Hex(),
			"username": user.Username,
			"email": user.Email,
		})
		c.Set("User", template.JS(url.QueryEscape(string(User))))
	}
	collection = db.C(ACCOUNTS)
	ac := Account{} // todo pool
	err = collection.FindId(bson.ObjectIdHex(user.Roles.Account.Hex())).One(&ac)
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
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func SetSettings (c *gin.Context) {
	account := getAccount(c)
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)
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
		response.Fail()
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()

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

	collection := db.C(ACCOUNTS)
	err = collection.UpdateId(account.ID, account)
	if err != nil {
		response.Fail()
		return
	}

	response.Finish()
}

func ChangePassword (c *gin.Context) {
	user := getUser(c)
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)
	var body struct {
		Confirm   string  `json:"confirm"`
		Password string  `json:"newPassword"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	// validate
	if len(body.Password) == 0 {
		response.ErrFor["newPassword"] = "required"
	}
	if len(body.Confirm) == 0 {
		response.ErrFor["confirm"] = "required"
	} else if body.Password != body.Confirm {
		response.Errors = append(response.Errors, "Passwords do not match.")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Errors = append(response.Errors, err.Error()) // TODO don't like that this error goes to client
		response.Fail()
		return
	}

	user.Password = string(hashedPassword)
	collection := db.C(USERS)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}

func ChangeIdentity (c *gin.Context) {
	user := getUser(c)
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)
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
		response.Fail()
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)

	{
		us := User{} // todo pool
		err = collection.Find(bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}).One(&us)
		if err != nil {
			response.Errors = append(response.Errors, "username or email already exist")
			response.Fail()
			return
		}
	}

	user.Username = body.Username
	user.Email = body.Email

	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}
	// TODO  patch admin and account
	response.Finish()
}

func providerSettings (c *gin.Context) {
	session := sessions.Default(c)
	session.Set("action", "/account/settings/")
	session.Save()
	startOAuth(c)
}

func settingsProvider (c *gin.Context, userGoth goth.User) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := &User{}
	err := collection.Find(bson.M{userGoth.Provider + ".id": userGoth.UserID}).One(user)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		session := sessions.Default(c)
		session.Set("oauthMessage", "Another user has already connected with that " + userGoth.Provider + " account")
		session.Save()
		c.Redirect(http.StatusFound, "/account/settings/")
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	user = getUser(c)

	user.updateProvider(userGoth)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		panic(err)
		return
	}

	c.Redirect(http.StatusFound, "/account/settings/")
}

func disconnectProvider (c *gin.Context) {
	user := getUser(c)
	user.disconnectProviderDB(c.Param("provider"))
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err := collection.UpdateId(user.ID, user)
	if err != nil {
		panic(err)
		return
	}

	c.Redirect(http.StatusFound, "/account/settings/")
}
