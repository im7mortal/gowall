package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
)

func renderStatuses(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.RegEx{
			Pattern: `^.*?` + name + `.*$`,
			Options: "i",
		}
	}

	pivot, ok := c.GetQuery("pivot")
	if ok && len(pivot) != 0 {
		query["pivot"] = bson.RegEx{
			Pattern: `^.*?` + pivot + `.*$`,
			Options: "i",
		}
	}

	var statuses []Status

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)

	Result := getData(c, collection.Find(query), &statuses)

	filters := Result["filters"].(gin.H)
	filters["name"] = name
	filters["pivot"] = pivot

	Results, err := json.Marshal(Result)
	if err != nil {
		EXCEPTION(err)
	}

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func createStatus(c *gin.Context) {
	response := newResponse(c)
	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not create statuses")
		response.Fail()
		return
	}

	status := Status{}

	err := json.NewDecoder(c.Request.Body).Decode(&status)
	if err != nil {
		EXCEPTION(err)
	}

	if len(status.Name) == 0 {
		response.ErrFor["name"] = "required"
	}

	if len(status.Pivot) == 0 {
		response.ErrFor["pivot"] = "required"
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	//duplicateStatusCheck
	_id := slugify(status.Pivot + " " + status.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)
	err = collection.FindId(_id).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That status+pivot is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// createStatus
	status.ID = _id
	err = collection.Insert(status)
	if err != nil {
		EXCEPTION(err)
	}
	response.Finish()
}

func readStatus(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)
	status := Status{}
	err := collection.FindId(c.Param("id")).One(&status)
		if err != nil {
		if err == mgo.ErrNotFound {
			renderStatus404(c)
			return
		}
		EXCEPTION(err)
	}

	json, err := json.Marshal(status)
	if err != nil {
		EXCEPTION(err.Error())
	}

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))
	c.HTML(http.StatusOK, "/admin/statuses/details/", c.Keys)
}

func updateStatus(c *gin.Context) {
	response := newResponse(c)
	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not create statuses.")
		response.Fail()
		return
	}

	status := Status{}

	err := json.NewDecoder(c.Request.Body).Decode(&status)
	if err != nil {
		EXCEPTION(err)
	}

	if len(status.Name) == 0 {
		response.ErrFor["name"] = "required"
	}

	if len(status.Pivot) == 0 {
		response.ErrFor["pivot"] = "required"
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	//duplicateStatusCheck
	_id := slugify(status.Pivot + " " + status.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)
	err = collection.FindId(_id).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That status+pivot is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		EXCEPTION(err)
	}

	// patchStatus
	status.ID = _id
	err = collection.RemoveId(c.Param("id")) // c.Param("id") is string/ no bson.ObjectID
	if err != nil {
		EXCEPTION(err)
	}
	err = collection.Insert(status)
	if err != nil {
		EXCEPTION(err)
	}

	response.Finish()
}

func deleteStatus(c *gin.Context) {
	response := newResponse(c)
	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not delete statuses.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)
	err := collection.RemoveId(c.Param("id")) // c.Param("id") is string/ no bson.ObjectID
	if err != nil {
		EXCEPTION(err)
	}

	response.Finish()
}
