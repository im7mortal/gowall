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

	res := gin.H{
		"_id":      user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"isActive": user.IsActive,
	}

	// TODO parrallel?

	if len(user.Roles.Admin.Hex()) != 0 {
		admin := Admin{}
		collection := db.C(ADMINS)
		collection.FindId(user.Roles.Admin).One(&admin)
		res["roles"] = gin.H{
			"admin": gin.H{
				"id_": admin.ID.Hex(),
				"name": gin.H{
					"full": admin.Name.Full,
				},
			},
		}
	}

	if len(user.Roles.Account.Hex()) != 0 {
		roles, ok := res["roles"].(gin.H)
		if !ok {
			roles = gin.H{}
		}
		account := Account{}
		collection := db.C(ACCOUNTS)
		collection.FindId(user.Roles.Account).One(&account)
		roles["account"] = gin.H{
			"id_": account.ID.Hex(),
			"name": gin.H{
				"full": account.Name.Full,
			},
		}
		res["roles"] = roles
	}

	json, err := json.Marshal(res)
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
	response.Init(c)
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
	_id := c.Param("id")
	err = collection.FindId(bson.ObjectIdHex(_id)).One(&user)
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
	err = collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{
				"username": body.Username,
			},
			bson.M{
				"email": body.Email,
			},
		},
		"_id": bson.M{ "$ne": bson.ObjectIdHex(_id)},
	}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "username or email already exist")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	user.Username = body.Username
	user.Email = body.Email

	// patchUser
	err = collection.UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"username": user.Username,
			"email": user.Email,
			"isActive": user.IsActive,
			"search": []string{user.Username, user.Email},
		},
	})
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}
	// patchAdmin
	collection = db.C(ADMINS)
	err = collection.Update(bson.M{ "user.id": user.ID},
		bson.M{ "$set": bson.M{ "user.name": user.Username},
	})

	// patchAccount
	collection = db.C(ACCOUNTS)
	err = collection.Update(bson.M{ "user.id": user.ID},
		bson.M{ "$set": bson.M{ "user.name": user.Username},
	})


	// populateRoles // TODO
	response.Finish()
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

func linkAdminToUser (c *gin.Context) {
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
		NewAdminId  string `json:"newAdminId"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
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
		panic(err)
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
		"_id": bson.M{ "$ne": bson.ObjectIdHex(userID)},
	}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "Another user is already linked to that admin.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	//patchUser
	//patchAdmin
	user := &User{}

	err = collection.FindId(bson.ObjectIdHex(userID)).One(user)

	if err != nil {
		if err != mgo.ErrNotFound {
		panic(err)
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
		"id_": userID,
		"timeCreated": userID, //TODO
		"username": user.Username,
		"search": []string{user.Username},
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

func unlinkAdminToUser (c *gin.Context) {
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
		panic(err)
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
		"id_": id_,
		"timeCreated": id_, //TODO
		"username": user.Username,
		"search": []string{},
		"roles": gin.H{
			"admin": nil,
		},
	}

	response.Finish()
}

func linkAccountToUser (c *gin.Context) {
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
		panic(err)
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
		panic(err)
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
		"_id": bson.M{ "$ne": bson.ObjectIdHex(userID)},
	}).One(nil)
	if err == nil {
		response.Errors = append(response.Errors, "Another user is already linked to that account.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	//patchUser
	//patchAccount
	user := &User{}

	err = collection.FindId(bson.ObjectIdHex(userID)).One(user)

	if err != nil {
		if err != mgo.ErrNotFound {
		panic(err)
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
		"id_": userID,
		"timeCreated": userID, //TODO
		"username": user.Username,
		"search": []string{user.Username},
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

func unlinkAccountToUser (c *gin.Context) {
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
		panic(err)
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
		"id_": id_,
		"timeCreated": id_, //TODO
		"username": user.Username,
		"search": []string{},
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
