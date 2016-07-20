package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"net/url"
)

type responseUser struct {
	Response
	User
}

func renderUsers(c *gin.Context) {
	query := bson.M{}

	username, ok := c.GetQuery("username")
	if ok && len(username) != 0 {
		query["username"] = bson.RegEx{
			Pattern: `^.*?` + username + `.*$`,
			Options: "i",
		}
	}

	isActive, ok := c.GetQuery("isActive")
	if ok && len(isActive) != 0 {
		query["isActive"] = isActive
	}

	roles, ok := c.GetQuery("roles")
	if ok && len(roles) != 0 {
		// roles.admin or roles.account
		query["roles."+roles] = bson.M{
			"$exists": true,
		}
	}

	type _user struct {
		ID       bson.ObjectId `bson:"_id" json:"_id"`
		Username string        `bson:"username" json:"username"`
		IsActive string        `bson:"isActive" json:"isActive"`
		Email    string        `bson:"email" json:"email"`
	}

	var results []_user

	db := getMongoDBInstance()
	defer db.Session.Close()

	collection := db.C(USERS)

	Result := getData(c, collection.Find(query), &results)

	filters := Result["filters"].(gin.H)
	filters["username"] = username
	filters["isActive"] = isActive
	filters["roles"] = roles

	Results, err := json.Marshal(Result)
	if err != nil {
		panic(err.Error())
	}

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func createUser(c *gin.Context) {
	response := responseUser{}
	response.Init(c)

	err := json.NewDecoder(c.Request.Body).Decode(&response.User)
	if err != nil {
		panic(err)
		return
	}

	// validate
	validateUsername(&response.User.Username, &response.Response)
	if response.HasErrors() {
		response.Fail()
		return
	}

	// duplicateUsernameCheck
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err = collection.Find(bson.M{"username": response.Username}).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That username is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// createUser
	response.User.ID = bson.NewObjectId()
	response.User.Search = []string{response.Username}

	err = collection.Insert(response.User)
	if err != nil {
		panic(err)
		return
	}
	response.Data["record"] = response.User
	response.Finish()
}

func readUser(c *gin.Context) {

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	json, err := json.Marshal(gin.H{
		"_id":      user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"isActive": user.IsActive,
	})
	if err != nil {
		panic(err.Error())
	}

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(url.QueryEscape(string(json))))
	c.HTML(http.StatusOK, "/admin/users/details/", c.Keys)
}

func changeDataUser(c *gin.Context) {
	response := Response{}

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		IsActive string `json:"isActive"`
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
	user := User{}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail()
		return
	}

	if len(user.IsActive) == 0 {
		user.IsActive = "no"
	}

	// duplicateUsernameCheck
	// duplicateEmailCheck
	{
		us := User{}
		err = collection.Find(bson.M{
			"$or": []bson.M{
				bson.M{
					"username": body.Username,
				},
				bson.M{
					"email": body.Email,
				},
			},
		}).One(&us)
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
	response.Success = true
	c.JSON(http.StatusOK, response)
}

func changePasswordUser(c *gin.Context) {
	response := Response{}
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

	// patchUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail()
		return
	}

	user.setPassword(body.Password)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}

func deleteUser(c *gin.Context) {
	admin := getAdmin(c)
	user := getUser(c)

	response := Response{}
	response.Init(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not delete users.")
		response.Fail()
		return
	}

	deleteID := c.Param("id")

	if deleteID == user.ID.Hex() {
		response.Errors = append(response.Errors, "You may not delete yourself from user.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err := collection.RemoveId(bson.ObjectIdHex(deleteID))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}
