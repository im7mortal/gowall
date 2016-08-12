package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"sync"
	"strings"
	"time"
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
		query["roles." + roles] = bson.M{
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
		EXCEPTION(err.Error())
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
		EXCEPTION(err)
		return
	}

	// validate
	if ok := rUsername.MatchString(strings.ToLower(response.User.Username)); !ok {
		response.Errors = append(response.Errors, `only use letters, numbers, -, _`)
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
		EXCEPTION(err)
	}

	// createUser
	response.User.ID = bson.NewObjectId()
	response.User.TimeCreated = time.Now()
	response.User.Search = []string{response.Username}

	err = collection.Insert(response.User)
	if err != nil {
		EXCEPTION(err)
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

	res := gin.H{
		"_id":      user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"isActive": user.IsActive,
	}

	res["roles"] = getRoles(db, &user)
	json, err := json.Marshal(res)
	if err != nil {
		EXCEPTION(err.Error())
	}

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))
	c.HTML(http.StatusOK, "/admin/users/details/", c.Keys)
}

func changeDataUser(c *gin.Context) {
	response := newResponse(c)

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail()
		return
	}
	err = user.changeIdentity(response)
	if err != nil {
		response.Fail()
		return
	}
	response.Finish()
}

func changePasswordUser(c *gin.Context) {
	response := newResponse(c)

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail()
		return
	}

	// patchUser
	err = user.changePassword(response)
	if err != nil {
		response.Fail()
		return
	}

	response.Finish()
}

func linkAdminToUser(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)
	//user := getUser(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not link users to admins.")
		response.Fail()
		return
	}
	var body struct {
		NewAdminId string `json:"newAdminId"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}
	if len(body.NewAdminId) == 0 {
		response.ErrFor["newAdminId"] = "required"
		response.Fail()
		return
	}

	// verifyAdmin
	db := getMongoDBInstance()
	defer db.Session.Close()

	collection := db.C(ADMINS)
	admin = &Admin{}
	err = collection.FindId(bson.ObjectIdHex(body.NewAdminId)).One(admin)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "Admin not found.")
		response.Fail()
		return
	}
	userID := c.Param("id")
	id_ := admin.User.ID.Hex()

	if len(id_) == 12 && id_ != userID {
		response.Errors = append(response.Errors, "Admin is already linked to a different user.")
		response.Fail()
		return
	}

	//duplicateLinkCheck
	collection = db.C(USERS)
	err = collection.Find(bson.M{
		"roles.admin": admin.ID,
		"_id":         bson.M{"$ne": bson.ObjectIdHex(userID)},
	}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "Another user is already linked to that admin.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	//patchUser
	//patchAdmin
	user := &User{}

	err = collection.FindId(bson.ObjectIdHex(userID)).One(user)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User not found.")
		response.Fail()
		return
	}
	err = admin.linkUser(db, user)
	if err != nil {
		response.Errors = append(response.Errors, "Something went wrong.")
		response.Fail()
		return
	}

	response.Data["user"] = gin.H{
		"id_":         userID,
		"timeCreated": user.TimeCreated.Format(ISOSTRING),
		"username":    user.Username,
		"search":      []string{user.Username},
		"roles": gin.H{
			"admin": gin.H{
				"id_": admin.ID.Hex(),
				"name": gin.H{
					"full": admin.Name.Full,
				},
			},
		},
	}
	response.Finish()
}

func unlinkAdminFromUser(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not link users to admins.")
		response.Fail()
		return
	}
	id_ := c.Param("id")
	if admin.ID.Hex() == id_ {
		response.Errors = append(response.Errors, "You may not unlink yourself from admin.")
		response.Fail()
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()

	collection := db.C(USERS)
	user := &User{}

	err := collection.FindId(bson.ObjectIdHex(id_)).One(user)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User not found.")
		response.Fail()
		return
	}
	admin = &Admin{}

	//patchUser
	//patchAdmin
	err = admin.unlinkUser(db, user)
	if err != nil {
		response.Errors = append(response.Errors, "Something went wrong.")
		response.Fail()
		return
	}
	response.Data["user"] = gin.H{
		"id_":         id_,
		"timeCreated": user.TimeCreated.Format(ISOSTRING),
		"username":    user.Username,
		"search":      []string{},
		"roles": gin.H{
			"admin": nil,
		},
	}

	response.Finish()
}

func linkAccountToUser(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not link users to admins.")
		response.Fail()
		return
	}
	var body struct {
		NewAccountId string `json:"newAccountId"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}
	if len(body.NewAccountId) == 0 {
		response.ErrFor["newAccountId"] = "required"
		response.Fail()
		return
	}

	// verifyAccount
	db := getMongoDBInstance()
	defer db.Session.Close()

	collection := db.C(ACCOUNTS)
	account := &Account{}
	err = collection.FindId(bson.ObjectIdHex(body.NewAccountId)).One(account)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "Account not found.")
		response.Fail()
		return
	}
	userID := c.Param("id")
	id_ := account.User.ID.Hex()

	if len(id_) == 12 && id_ != userID {
		response.Errors = append(response.Errors, "Account is already linked to a different user.")
		response.Fail()
		return
	}

	//duplicateLinkCheck
	collection = db.C(USERS)
	err = collection.Find(bson.M{
		"roles.account": account.ID,
		"_id":           bson.M{"$ne": bson.ObjectIdHex(userID)},
	}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "Another user is already linked to that account.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	//patchUser
	//patchAccount
	user := &User{}

	err = collection.FindId(bson.ObjectIdHex(userID)).One(user)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User not found.")
		response.Fail()
		return
	}
	err = account.linkUser(db, user)
	if err != nil {
		response.Errors = append(response.Errors, "Something went wrong.")
		response.Fail()
		return
	}

	response.Data["user"] = gin.H{
		"id_":         userID,
		"timeCreated": user.TimeCreated.Format(ISOSTRING),
		"username":    user.Username,
		"search":      []string{user.Username},
		"roles": gin.H{
			"account": gin.H{
				"id_": account.ID.Hex(),
				"name": gin.H{
					"full": account.Name.Full,
				},
			},
		},
	}
	response.Finish()
}

func unlinkAccountFromUser(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not link users to accounts.")
		response.Fail()
		return
	}
	id_ := c.Param("id")

	db := getMongoDBInstance()
	defer db.Session.Close()

	collection := db.C(USERS)
	user := &User{}

	err := collection.FindId(bson.ObjectIdHex(id_)).One(user)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User not found.")
		response.Fail()
		return
	}
	account := &Account{}

	//patchUser
	//patchAccount
	err = account.unlinkUser(db, user)
	if err != nil {
		response.Errors = append(response.Errors, "Something went wrong.")
		response.Fail()
		return
	}
	response.Data["user"] = gin.H{
		"id_":         id_,
		"timeCreated": user.TimeCreated.Format(ISOSTRING),
		"username":    user.Username,
		"search":      []string{},
		"roles": gin.H{
			"account": nil,
		},
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

func getRoles(db *mgo.Database, user *User) (roles gin.H) {
	var wg sync.WaitGroup
	var err error
	roles = gin.H{}
	if len(user.Roles.Admin.Hex()) != 0 {
		wg.Add(1)
		go func() {
			admin := Admin{}
			err = db.C(ADMINS).FindId(user.Roles.Admin).One(&admin)
			if err != nil {
				if err != mgo.ErrNotFound {
					EXCEPTION(err)
				}
			} else {
				roles["admin"] = gin.H{
					"id_": admin.ID.Hex(),
					"name": gin.H{
						"full": admin.Name.Full,
					},
				}
			}
			wg.Done()
		}()
	}

	if len(user.Roles.Account.Hex()) != 0 {
		wg.Add(1)
		go func() {
			account := Account{}
			err = db.C(ACCOUNTS).FindId(user.Roles.Account).One(&account)
			if err != nil {
				if err != mgo.ErrNotFound {
					EXCEPTION(err)
				}
			} else {
				roles["account"] = gin.H{
					"id_": account.ID.Hex(),
					"name": gin.H{
						"full": account.Name.Full,
					},
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return
}
