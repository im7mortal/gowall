package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type responseAccount struct {
	Response
	Account
}

func renderAccounts(c *gin.Context) {
	query := bson.M{}

	search, ok := c.GetQuery("search")
	if ok && len(search) != 0 {
		query["search"] = bson.RegEx{
			Pattern: `^.*?` + search + `.*$`,
			Options: "i",
		}
	}

	status, ok := c.GetQuery("status")
	if ok && len(status) != 0 {
		query["status"] = status
	}

	var results []Account

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)

	Result := getData(c, collection.Find(query), &results)

	filters := Result["filters"].(gin.H)
	filters["search"] = search
	filters["status"] = status

	Results, err := json.Marshal(Result)
	if err != nil {
		EXCEPTION(err)
	}

	if XHR(c) {
		handleXHR(c, Results)
		return
	}
	c.Set("Results", template.JS(getEscapedString(string(Results))))

	var statuses []Status
	collection = db.C(STATUSES)
	err = collection.Find(nil).All(&statuses)

	// preparing for js.  Don't like it.
	// https://groups.google.com/forum/#!topic/golang-nuts/0HJoROz2TMo
	// https://play.golang.org/p/M_AoMQwtFt
	// 10 july 2016 wasn't expected
	var statusesS string = `<option value="">-- any --</option>`
	for _, status := range statuses {
		statusesS += `<option value="` + status.ID + `">` + status.Name + `</option>`
	}
	c.Set("Statuses", template.JS(statusesS))
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func createAccount(c *gin.Context) {
	response := Response{}
	response.Init(c)

	var body struct {
		Name string `json:"name.full"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}
	if len(body.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
		response.Fail()
		return
	}

	account := Account{}

	// handleName
	account.Name.Full = slugifyName(body.Name)

	// duplicateAccount
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	err = collection.Find(bson.M{
		"name.full": account.Name.Full,
	}).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That account already exists.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// handleName
	name := strings.Split(account.Name.Full, " ")
	account.Name.First = name[0]
	if len(name) > 1 {
		if len(name) == 2 {
			account.Name.Last = name[1]
			account.Name.Middle = ""
		}
		if len(name) == 3 {
			account.Name.Middle = name[2]
		}
	}

	account.Search = []string{account.Name.First, account.Name.Middle, account.Name.Last}

	// createAccount
	account.ID = bson.NewObjectId()
	admin := getAdmin(c)
	account.UserCreated.ID = admin.ID
	account.UserCreated.Name = admin.Name.Full
	account.UserCreated.Time = time.Now()
	err = collection.Insert(account)
	if err != nil {
		EXCEPTION(err)
	}
	response.Data["record"] = account
	response.Finish()
}

func readAccount(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	account := Account{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&account)
	if err != nil {
		if err == mgo.ErrNotFound {
			renderStatus404(c)
			return
		}
		EXCEPTION(err)
	}
	json, err := json.Marshal(account)
	if err != nil {
		EXCEPTION(err)
	}
	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))

	var statuses []Status
	collection = db.C(STATUSES)
	err = collection.Find(nil).All(&statuses)

	// preparing for js.  Don't like it.
	// https://groups.google.com/forum/#!topic/golang-nuts/0HJoROz2TMo
	// https://play.golang.org/p/M_AoMQwtFt
	// 10 july 2016 wasn't expected
	var statusesS string = `<option value="">-- any --</option>`
	for _, status := range statuses {
		statusesS += `<option value="` + status.ID + `">` + status.Name + `</option>`
	}
	c.Set("Statuses", template.JS(statusesS))
	c.HTML(http.StatusOK, "/admin/accounts/details/", c.Keys)
}

func updateAccount(c *gin.Context) {
	response := newResponse(c)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	account := &Account{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(account)
	if err != nil {
		EXCEPTION(err)
	}
	err = account.changeData(response)
	if err != nil {
		response.Fail()
		return
	}
	response.Data["account"] = account
	response.Finish()
}

func linkUserToAccount(c *gin.Context) {
	response := responseAccount{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not link accounts to users.")
		response.Fail()
		return
	}

	var req struct {
		NewUsername string `json:"newUsername"`
	}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		EXCEPTION(err)
	}

	if len(req.NewUsername) == 0 {
		response.ErrFor["newUsername"] = "required"
		response.Errors = append(response.Errors, "required")
		response.Fail()
		return
	}

	//verifyUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := &User{}
	err = collection.Find(bson.M{"username": req.NewUsername}).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User not found.")
		response.Fail()
		return
	}
	id := c.Param("id")
	if user.Roles.Account.Hex() == id {
		response.Errors = append(response.Errors, "User is already linked to a different account.")
		response.Fail()
		return
	}

	account := Account{}
	// duplicateLinkCheck
	collection = db.C(ACCOUNTS)
	err = collection.Find(
		bson.M{
			"user.id": id,
			"_id": bson.M{
				"user.id": id,
			},
		}).One(&account) // reuse account. If it will be used it mean that user already linked.

	if err == nil {
		response.Errors = append(response.Errors, "Another account is already linked to that user.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	account.ID = bson.ObjectIdHex(id)
	// patchUser patchAccount
	err = account.linkUser(db, user)

	if err != nil {

	}

	err = collection.FindId(bson.ObjectIdHex(id)).One(&response.Account)

	if err != nil {
		EXCEPTION(err)
	}

	response.Data["account"] = response.Account
	response.Finish()
}

func unlinkUserFromAccount(c *gin.Context) {
	response := responseAccount{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not unlink accounts to users.")
		response.Fail()
		return
	}
	id_ := c.Param("id")
	response.ErrFor = map[string]string{} // in that handler it required (non standard behavior from node)

	// patchUser
	db := getMongoDBInstance()
	defer db.Session.Close()

	collection := db.C(ACCOUNTS)
	account := &Account{}

	err := collection.FindId(bson.ObjectIdHex(id_)).One(account)

	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
	}

	collection = db.C(USERS)
	user := &User{}

	err = collection.FindId(account.User.ID).One(user)

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
		EXCEPTION(err)
	}

	response.Data["account"] = response.Account
	response.Finish()
}

func newNote(c *gin.Context) {
	user := getUser(c)
	response := responseAccount{}
	response.Init(c)

	// validate
	var body struct {
		Data string `json:"data"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}
	if len(body.Data) == 0 {
		response.Errors = append(response.Errors, "Data is required.")
		response.Fail()
		return
	}

	// addNote
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	account := &Account{}
	err = collection.UpdateId(bson.ObjectIdHex(c.Param("id")),
		bson.M{"$push": bson.M{"notes": bson.M{
			"_id":  bson.NewObjectId(),
			"data": body.Data,
			"userCreated": bson.M{
				"id":   user.ID,
				"name": user.Username,
				"time": time.Now(),
			},
		}},
		})
	if err != nil {
		EXCEPTION(err)
	}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(account)
	if err != nil {
		EXCEPTION(err)
	}
	response.Data["account"] = account
	response.Finish()
}

func newStatus(c *gin.Context) {
	user := getUser(c)
	response := responseAccount{}
	response.Init(c)

	// validate
	var body struct {
		StatusID string `json:"id"`
		Name     string `json:"name"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)
	if err != nil {
		EXCEPTION(err)
	}
	if len(body.StatusID) == 0 {
		response.Errors = append(response.Errors, "Please choose a status.")
		response.Fail()
		return
	}

	// addStatus
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	account := &Account{}
	statusToAdd := bson.M{
		"_id":  body.StatusID,
		"name": body.Name,
		"userCreated": bson.M{
			"id":   user.ID,
			"name": user.Username,
			"time": time.Now(),
		},
	}
	err = collection.UpdateId(bson.ObjectIdHex(c.Param("id")),
		bson.M{
			"$push": bson.M{"statusLog": statusToAdd},
			"$set":  bson.M{"status": statusToAdd},
		})
	if err != nil {
		EXCEPTION(err)
	}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(account)
	if err != nil {
		EXCEPTION(err)
	}
	response.Data["account"] = account
	response.Finish()
}

func deleteAccount(c *gin.Context) {
	response := Response{}
	response.Init(c)

	// validate
	if ok := getAdmin(c).IsMemberOf(ROOTGROUP); !ok {
		response.Errors = append(response.Errors, "You may not delete accounts.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ACCOUNTS)
	err := collection.RemoveId(bson.ObjectIdHex(c.Param("id")))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}
