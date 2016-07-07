package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"gopkg.in/mgo.v2"
	"net/url"
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

	var results []Status

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)

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

func createStatus(c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.BindContext(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not create statuses")
		response.Fail()
		return
	}

	status := Status{}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&status)
	if err != nil {
		panic(err)
		return
	}
	// clean errors from client
	response.CleanErrors()

	if len(status.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(status.Pivot) == 0 {
		response.Errors = append(response.Errors, "A pivot is required")
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
	status_ := Status{}
	err = collection.Find(bson.M{"_id": _id}).One(&status_)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That status+pivot is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// createStatus
	status.ID = _id

	err = collection.Insert(status)
	if err != nil {
		panic(err)
		return
	}
	response.Success = true
	c.JSON(http.StatusOK, response)
}

func readStatus(c *gin.Context) {

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)
	status := Status{}
	collection.Find(bson.M{"_id": c.Param("id")}).One(&status)
	json, _ := json.Marshal(status)

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(url.QueryEscape(string(json))))
	c.HTML(http.StatusOK, "/admin/statuses/details/", c.Keys)
}

func updateStatus(c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.BindContext(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not create statuses")
		response.Fail()
		return
	}

	status := Status{}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&status)
	if err != nil {
		panic(err)
		return
	}
	// clean errors from client
	response.CleanErrors()

	if len(status.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(status.Pivot) == 0 {
		response.Errors = append(response.Errors, "A pivot is required")
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
	status_ := Status{}
	err = collection.FindId(_id).One(&status_)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That status+pivot is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// patchStatus
	status.ID = _id
	err = collection.RemoveId(c.Param("id"))
	//println(err.Error())
	err = collection.Insert(status)
	//println(err.Error())
	if err != nil {
		panic(err)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

func deleteStatus(c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.BindContext(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not delete statuses.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(STATUSES)
	err := collection.RemoveId(c.Param("id"))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}
