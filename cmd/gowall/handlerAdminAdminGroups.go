package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

type responseAdminGroup struct {
	Response
	AdminGroup
}

func renderAdminGroups(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.RegEx{
			Pattern: `^.*?` + name + `.*$`,
			Options: "i",
		}
	}

	var results []AdminGroup

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)

	Result := getData(c, collection.Find(query), &results)

	filters := Result["filters"].(gin.H)
	filters["name"] = name

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

func createAdminGroup(c *gin.Context) {
	response := responseAdminGroup{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not create statuses")
		response.Fail()
		return
	}

	response.AdminGroup.DecodeRequest(c)

	if len(response.AdminGroup.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
		response.Fail()
		return
	}

	//duplicateAdminGroupCheck
	response.AdminGroup.Name = slugifyName(response.AdminGroup.Name)
	response.AdminGroup.ID = slugify(response.AdminGroup.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	err := collection.FindId(response.ID).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That group already exists.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// createAdminGroup
	err = collection.Insert(response.AdminGroup)
	if err != nil {
		EXCEPTION(err)
	}
	response.Finish()
}

func readAdminGroup(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	adminGroup := AdminGroup{}
	err := collection.FindId(c.Param("id")).One(&adminGroup)
	if err != nil {
		if err == mgo.ErrNotFound {
			renderStatus404(c)
			return
		}
		EXCEPTION(err)
	}
	json, err := json.Marshal(adminGroup)
	if err != nil {
		EXCEPTION(err)
	}
	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))
	c.HTML(http.StatusOK, "/admin/admin-groups/details/", c.Keys)
}

func updateAdminGroup(c *gin.Context) {
	response := responseAdminGroup{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not update admin groups.")
		response.Fail()
		return
	}

	response.AdminGroup.DecodeRequest(c)

	if len(response.AdminGroup.Name) == 0 {
		response.ErrFor["name"] = "required"
		response.Fail()
		return
	}

	//duplicateAdminGroupCheck
	response.AdminGroup.ID = slugify(response.AdminGroup.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	adm := AdminGroup{}
	err := collection.FindId(response.AdminGroup.ID).One(&adm)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That admin group is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// patchAdminGroup
	// _id is slugified name so first delete second insert
	err = collection.RemoveId(c.Param("id"))
	if err != nil {
		EXCEPTION(err)
	}
	err = collection.Insert(response.AdminGroup)
	if err != nil {
		EXCEPTION(err)
	}

	response.Finish()
}

func updatePermissionsAdminGroup(c *gin.Context) {
	response := responseAdminGroup{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not change the permissions of admin groups.")
		response.Fail()
		return
	}

	response.AdminGroup.DecodeRequest(c)

	if len(response.AdminGroup.Permissions) == 0 {
		response.ErrFor["permissions"] = "required"
		response.Fail()
		return
	}

	//patchAdminGroup
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)

	err := collection.UpdateId(c.Param("id"), response.AdminGroup) // id is string
	if err != nil {
		EXCEPTION(err)
	}

	response.Finish()
}

func deleteAdminGroup(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not delete admin groups.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	err := collection.RemoveId(c.Param("id"))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}
