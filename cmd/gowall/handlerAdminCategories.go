package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
)

func AdminCategoriesRender(c *gin.Context) {
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

	type _category struct {
		ID bson.ObjectId `bson:"_id" json:"_id"`
		Name string `bson:"name" json:"name"`
		Pivot string `bson:"pivot" json:"pivot"`
	}

	var results []_category

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)

	Result := getData(c, collection.Find(query), &results)

	Results, _ := json.Marshal(Result)

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))

	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func CategoryRender(c *gin.Context) {

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	collection.FindId(c.Param("id")).One(&user)
	userJSON, _ := json.Marshal(user)
	c.Set("Record", string(userJSON))
	render, _ := TemplateStorage["/admin/users/details/"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}
