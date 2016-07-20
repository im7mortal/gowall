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

func renderCategories(c *gin.Context) {
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

	var results []Category

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)

	Result := getData(c, collection.Find(query), &results)

	filters := Result["filters"].(gin.H)
	filters["name"] = name
	filters["pivot"] = pivot

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

func createCategory(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not create categories")
		response.Fail()
		return
	}

	category := Category{}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&category)
	if err != nil {
		panic(err)
		return
	}

	if len(category.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(category.Pivot) == 0 {
		response.Errors = append(response.Errors, "A pivot is required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	//duplicateCategoryCheck
	_id := slugify(category.Pivot + " " + category.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)
	err = collection.FindId(_id).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That category+pivot is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// createCategory
	category.ID = _id

	err = collection.Insert(category)
	if err != nil {
		panic(err)
		return
	}
	response.Finish()
}

func updateCategory(c *gin.Context) {
	response := Response{}
	response.Init(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not create categories")
		response.Fail()
		return
	}

	category := Category{}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&category)
	if err != nil {
		panic(err)
		return
	}

	if len(category.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(category.Pivot) == 0 {
		response.Errors = append(response.Errors, "A pivot is required")
	}

	if response.HasErrors() {
		response.Fail()
		return
	}

	//duplicateCategoryCheck
	_id := slugify(category.Pivot + " " + category.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)
	err = collection.FindId(_id).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That category+pivot is already taken.")
		response.Fail()
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// patchCategory
	category.ID = _id
	err = collection.RemoveId(c.Param("id"))
	err = collection.Insert(category)
	if err != nil {
		panic(err)
		return
	}

	response.Finish()
}

func renderCategory(c *gin.Context) {

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)
	category := Category{}
	err := collection.FindId(c.Param("id")).One(&category)
	if err != nil {
		if err == mgo.ErrNotFound {
			Status404Render(c)
			return
		}
		panic(err)
	}

	json, _ := json.Marshal(category)

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(url.QueryEscape(string(json))))
	c.HTML(http.StatusOK, "/admin/categories/details/", c.Keys)
}

func deleteCategory(c *gin.Context) {
	admin := getAdmin(c)

	response := Response{}
	response.Init(c)

	// validate
	ok := admin.IsMemberOf(ROOTGROUP)
	if !ok {
		response.Errors = append(response.Errors, "You may not delete categories.")
		response.Fail()
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)
	err := collection.RemoveId(c.Param("id"))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail()
		return
	}

	response.Finish()
}
