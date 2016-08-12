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

type responseAdmin struct {
	Response
	Admin
}

func renderAdmins(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.RegEx{
			Pattern: `^.*?` + name + `.*$`,
			Options: "i",
		}
	}

	var results []Admin

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	Result := getData(c, collection.Find(query), &results)

	Results, err := json.Marshal(Result)
	if err != nil {
		EXCEPTION(err)
	}
	if XHR(c) {
		handleXHR(c, Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func createAdmin(c *gin.Context) {
	response := responseAdmin{}
	response.Init(c)
	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not create administrators")
		response.Fail()
		return
	}

	response.Admin.DecodeRequest(c)
	if len(response.Name.Full) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	// handleName
	response.Name.Full = slugifyName(response.Name.Full)

	// duplicateAdministrator
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	err := collection.Find(bson.M{"name.full": response.Name.Full}).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That administrator already exists.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// handleName
	name := strings.Split(response.Name.Full, " ")
	response.Name.First = name[0]
	if len(name) == 2 {
		response.Name.Last = name[1]
		response.Admin.Name.Middle = ""
	}
	if len(name) == 3 {
		response.Name.Middle = name[2]
	}
	response.Admin.Search = []string{response.Name.First, response.Name.Middle, response.Name.Last}
	response.Admin.Permissions = []Permission{}
	response.Admin.Groups = []string{}

	// createAdministrator
	response.Admin.ID = bson.NewObjectId()
	response.Admin.TimeCreated = time.Now()
	err = collection.Insert(response.Admin)
	if err != nil {
		EXCEPTION(err)
	}
	response.Data["record"] = response
	response.Finish()
}

func readAdmin(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	admin := Admin{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&admin)
	if err != nil {
		if err == mgo.ErrNotFound {
			renderStatus404(c)
			return
		}
		EXCEPTION(err)
	}

	// populateGroups
	collection = db.C(ADMINGROUPS)
	adminGroups, err := admin.populateGroups(db)
	if len(admin.Permissions) == 0 {
		admin.Permissions = []Permission{}
	}
	if err != nil {
		// mgo.ErrNotFound is not possible. "Root" group must be.
		EXCEPTION(err)
	}

	results, err := json.Marshal(admin)
	if err != nil {
		EXCEPTION(err)
	}
	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", results)
		return
	}

	// preparing for js.  Don't like it.
	// https://groups.google.com/forum/#!topic/golang-nuts/0HJoROz2TMo
	// https://play.golang.org/p/M_AoMQwtFt
	// 10 july 2016 wasn't expected
	// !!!  somehow label broke it too
	var adminGroupsS string
	for _, adminGroup := range adminGroups {
		adminGroupsS += `<option value="` + adminGroup.ID + `">` + adminGroup.Name + `</option>`
	}
	c.Set("Groups", template.JS(adminGroupsS))
	c.Set("Record", template.JS(getEscapedString(string(results))))
	c.HTML(http.StatusOK, "/admin/administrators/details/", c.Keys)
}

func updateAdmin(c *gin.Context) {
	response := responseAdmin{}
	response.Init(c)

	err := json.NewDecoder(c.Request.Body).Decode(&response.Admin.Name)
	if err != nil {
		EXCEPTION(err)
	}

	if len(response.Name.First) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(response.Name.Last) == 0 {
		response.Errors = append(response.Errors, "A lastname is required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	response.Admin.Name.Full = response.Admin.Name.First + " " + response.Admin.Name.Last

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	// patchAdministrator
	err = collection.UpdateId(bson.ObjectIdHex(c.Param("id")), bson.M{
		"$set": bson.M{
			"name": response.Admin.Name,
		},
	})
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}

func updatePermissionsAdmin(c *gin.Context) {
	response := responseAdmin{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not change the permissions of admins.")
		response.Fail()
		return
	}

	response.Admin.DecodeRequest(c)
	if len(response.Permissions) == 0 {
		response.ErrFor["permissions"] = "required"
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	for _, permission := range response.Admin.Permissions {
		permission.ID = bson.NewObjectId() // we lost ID every time. it can be solved easy
	}

	//patchAdmin
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	id := bson.ObjectIdHex(c.Param("id"))
	err := collection.UpdateId(id, bson.M{
		"$set": bson.M{
			"permissions": response.Admin.Permissions,
		},
	})
	if err != nil {
		EXCEPTION(err)
	}

	admin = &Admin{}
	err = collection.FindId(id).One(admin)
	if err != nil {
		EXCEPTION(err)
	}
	admin.populateGroups(db)
	response.Data["admin"] = admin
	response.Finish()
}

func updateGroupsAdmin(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not change the group memberships of admins.")
		response.Fail()
		return
	}

	var req struct {
		Groups []AdminGroup `json:"groups"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		EXCEPTION(err)
	}

	if len(req.Groups) == 0 {
		response.ErrFor["groups"] = "required"
		response.Fail()
		return
	}

	setGroups := []string{}
	for _, group := range req.Groups {
		setGroups = append(setGroups, group.ID)
	}

	//patchAdmin
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	id := bson.ObjectIdHex(c.Param("id"))
	err = collection.UpdateId(id, bson.M{
		"$set": bson.M{
			"groups": setGroups,
		},
	})
	if err != nil {
		EXCEPTION(err)
	}
	admin = &Admin{}
	err = collection.FindId(id).One(admin)
	if err != nil {
		EXCEPTION(err)
	}

	admin.populateGroups(db)
	response.Data["admin"] = admin
	response.Finish()
}

func linkUserToAdmin(c *gin.Context) {
	response := responseAdmin{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not link admins to users.")
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
	if user.Roles.Admin.Hex() == id {
		response.Errors = append(response.Errors, "User is already linked to a different admin.")
		response.Fail()
		return
	}

	// duplicateLinkCheck
	collection = db.C(ADMINS)
	err = collection.Find(
		bson.M{
			"user.id": id,
			"_id": bson.M{
				"user.id": id,
			},
		}).One(&admin) // reuse admin. If it will be used it mean that user already linked.

	if err == nil {
		response.Errors = append(response.Errors, "Another admin is already linked to that user.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	admin.ID = bson.ObjectIdHex(id)
	// patchUser patchAdministrator
	err = admin.linkUser(db, user)

	if err != nil {
		EXCEPTION(err)
	}

	err = collection.FindId(bson.ObjectIdHex(id)).One(&response.Admin)

	if err != nil {
		EXCEPTION(err)
	}

	response.Data["admin"] = response.Admin
	response.Finish()
}

func unlinkUserFromAdmin(c *gin.Context) {
	response := responseAdmin{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not change the permissions of admin groups.")
		response.Fail()
		return
	}
	id := c.Param("id")
	if admin.ID.Hex() == id {
		response.Errors = append(response.Errors, "You may not unlink yourself from admin.")
		response.Fail()
		return
	}
	response.ErrFor = map[string]string{}

	// patchUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err := collection.Update(bson.M{"roles.admin": bson.ObjectIdHex(id)}, bson.M{
		"$set": bson.M{"roles.admin": ""},
	})
	if err != nil {
		if err != mgo.ErrNotFound {
			EXCEPTION(err)
		}
		response.Errors = append(response.Errors, "User not found.")
		response.Fail()
		return
	}

	// patchAdministrator
	collection = db.C(ADMINS)
	err = collection.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{"user": bson.M{}},
	})

	if err != nil {
		EXCEPTION(err)
	}

	response.Data["admin"] = response.Admin
	response.Finish()
}

func deleteAdmin(c *gin.Context) {
	response := Response{}
	response.Init(c)

	// validate
	if ok := getAdmin(c).IsMemberOf(ROOTGROUP); !ok {
		response.Errors = append(response.Errors, "You may not delete administrators.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	err := collection.RemoveId(bson.ObjectIdHex(c.Param("id")))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}
