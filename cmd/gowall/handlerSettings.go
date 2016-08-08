package main

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

func renderAccountSettings(c *gin.Context) {
	sess := sessions.Default(c)

	user := getUser(c)
	injectSocials(c)
	doUserHasSocials(c, user)

	public := sess.Get("public")
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user = &User{}
	err := collection.FindId(bson.ObjectIdHex(public.(string))).One(user)
	if err != nil {
		panic(err)
	}
	if len(user.Username) != 0 {
		User, _ := json.Marshal(gin.H{
			"_id":      user.ID.Hex(),
			"username": user.Username,
			"email":    user.Email,
		})
		c.Set("User", template.JS(getEscapedString(string(User))))
	}
	collection = db.C(ACCOUNTS)
	ac := Account{}
	err = collection.FindId(user.Roles.Account).One(&ac)
	if err != nil {
		panic(err)
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
		c.Set("Account", template.JS(getEscapedString(string(Account))))
	}
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func setSettings(c *gin.Context) {
	account := getAccount(c)
	response := Response{}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)
	var body struct {
		First   string `json:"first"`
		Middle  string `json:"middle"`
		Last    string `json:"last"`
		Company string `json:"company"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
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

func changePassword(c *gin.Context) {
	user := getUser(c)
	response := Response{}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)
	var body struct {
		Confirm  string `json:"confirm"`
		Password string `json:"newPassword"`
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

	user.setPassword(body.Password)

	collection := db.C(USERS)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}

func changeIdentity(c *gin.Context) {
	user := getUser(c)
	response := Response{}
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	response.Init(c)
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	validateUsername(&body.Username, &response)
	validateEmail(&body.Email, &response)

	if response.HasErrors() {
		response.Fail()
		return
	}
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)

	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{"username": body.Username},
			bson.M{"email": body.Email},
		},
	}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "That username or email already exist.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	user.Username = body.Username
	user.Email = body.Email

	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	err = updateRoles(db, user)

	if err != nil {
		panic(err)
	}
	response.Finish()
}

func settingsProvider_(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("action", "/account/settings/")
	session.Save()
	startOAuth(c)
}

func settingsProvider(c *gin.Context, userGoth goth.User) {
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
	}

	c.Redirect(http.StatusFound, "/account/settings/")
}

func disconnectProvider(c *gin.Context) {
	user := getUser(c)
	user.disconnectProviderDB(c.Param("provider"))
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err := collection.UpdateId(user.ID, user)
	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusFound, "/account/settings/")
}
