package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"gopkg.in/mgo.v2"
	"net/url"
	"strings"
)

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

	Results, _ := json.Marshal(Result)

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

type responseAdminGroup struct {
	Response
	AdminGroup
}

func createAdminGroup(c *gin.Context) {
	response := responseAdminGroup{}
	defer response.Recover(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not create statuses")
		response.Fail(c)
		return
	}

	response.AdminGroup.DecodeRequest(c)

	if len(response.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	//duplicateAdminGroupCheck
	response.Name = slugifyName(response.Name)
	response.ID = slugify(response.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	err := collection.FindId(response.ID).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That group already exists.")
		response.Fail(c)
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	// createAdminGroup
	err = collection.Insert(response.AdminGroup) // todo I think mgo's behavior isn't expected
	if err != nil {
		panic(err)
		return
	}
	response.Success = true
	c.JSON(http.StatusOK, response)
}

func readAdminGroup(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	adminGroup := AdminGroup{}
	err := collection.FindId(c.Param("id")).One(&adminGroup)
	if err != nil {
		if err == mgo.ErrNotFound {
			Status404Render(c)
			return
		}
		panic(err)
	}
	json, err := json.Marshal(adminGroup)
	if err != nil {
		panic(err)
	}
	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))
	c.HTML(http.StatusOK, "/admin/admin-groups/details/", c.Keys)
}

/**
	TODO can be problem
 */
func getEscapedString(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func updateAdminGroup(c *gin.Context) {
	response := responseAdminGroup{}
	defer response.Recover(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not update admin groups.")
		response.Fail(c)
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&response)
	if err != nil {
		panic(err)
	}
	// clean errors from client
	response.CleanErrors()

	if len(response.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	//duplicateAdminGroupCheck
	response.ID = slugify(response.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	err = collection.FindId(response.ID).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That admin group is already taken.")
		response.Fail(c)
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// patchAdminGroup
	err = collection.RemoveId(c.Param("id"))
	if err != nil {
		panic(err)
	}
	err = collection.Insert(response.AdminGroup)
	if err != nil {
		panic(err)
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

func updateAdminGroupPermissions(c *gin.Context) {
	response := responseAdminGroup{}
	defer response.Recover(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not change the permissions of admin groups.")
		response.Fail(c)
		return
	}

	response.AdminGroup.DecodeRequest(c)

	if len(response.Permissions) == 0 {
		response.Errors = append(response.Errors, "required")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	//patchAdminGroup
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)

	err := collection.UpdateId(c.Param("id"), response.AdminGroup)
	if err != nil {
		panic(err)
	}

	response.Finish(c)
}

func deleteAdminGroup(c *gin.Context) {
	response := Response{} // todo sync.Pool
	defer response.Recover(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not delete admin groups.")
		response.Fail(c)
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	err := collection.RemoveId(c.Param("id"))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}
