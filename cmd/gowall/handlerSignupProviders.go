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

	c.Set("emailExist", len(userGoth.Email) != 0)
	render, _ := TemplateStorage["/signup/social/"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func SignUpSocial(c *gin.Context) {

	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)

	var body struct {
		Email   string  `json:"email"`
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
	session := sessions.Default(c)

	socialProfile_, ok := session.Get("socialProfile").(string)
	if ok && len(socialProfile_) == 0 {
		return
	}
	socialProfile := goth.User{}
	err = json.Unmarshal([]byte(socialProfile_), &socialProfile)
	if err != nil {

	}

	var username string
	if len(socialProfile.Name) != 0 {
		username = socialProfile.Name
	} else if len(socialProfile.UserID) != 0 {
		username = socialProfile.UserID
	}
	reg, err := regexp.Compile(`/[^a-zA-Z0-9\-\_]/g`)


	usernameSrc := []byte(username)
	reg.ReplaceAll(usernameSrc, []byte(""))
	username = string(usernameSrc)


	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	us := User{} // todo pool
	err = collection.Find(bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}).One(&us)

	if len(us.Username) != 0 {
		response.Fail(c)
		return
	}
	if len(us.Email) != 0 {
		response.Fail(c)
		return
	}

	us.ID = bson.NewObjectId()
	us.IsActive = "yes"
	us.Username = us.ID.Hex()
	us.Email = strings.ToLower(email)
	us.Search = []string{username, email}

	//provider, _ := session.Get("provider").(string)

	us.Facebook = vendorOauth{}
	us.Facebook.ID = socialProfile.UserID

	err = collection.Insert(us)

	if err != nil {
		println(err.Error())
		//todo error handling
		return
	}
	ac := Account{}

	ac.ID = bson.NewObjectId()

	us.Roles.Account = ac.ID

	err = collection.UpdateId(us.ID, us)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	if config.RequireAccountVerification {
		ac.IsVerified = "no"
	} else {
		ac.IsVerified = "yes"
	}
	ac.Name.Full = username
	ac.User.ID = us.ID
	ac.User.Name = us.Username
	ac.Search = []string{username}

	collection = db.C(ACCOUNTS)
	err = collection.Insert(ac)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}
	sess := sessions.Default(c)
	sess.Set("public", us.ID.Hex())
	sess.Save()

	response.Success = true
	c.JSON(http.StatusOK, response)
}