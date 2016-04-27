package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"html/template"
)

func AdminCategoriesRender(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.M{
			"$regex": "/^.*?" + name + ".*$/i",
		}
	}

	pivot, ok := c.GetQuery("pivot")
	if ok && len(pivot) != 0 {
		query["pivot"] = bson.M{
			"$regex": "/^.*?" + pivot + ".*$/i",
		}
	}

	limit_ := c.DefaultQuery("limit", "20")
	limit, _ := strconv.ParseInt(limit_, 0, 0)
	if limit > 100 {
		limit = 100
	}

	page_ := c.DefaultQuery("page", "0")
	page, _ := strconv.ParseInt(page_, 0, 0)

	sort := c.DefaultQuery("sort", "_id")

	var results []Category

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)
	collection.Find(query).Skip(int(limit * page)).Sort(sort).Limit(int(limit)).All(&results)

	categoies := []gin.H{}

	for _, adminGroup := range results {
		categoies = append(categoies, gin.H{
			"_id": adminGroup.ID.Hex(),
			"name": adminGroup.Name,
			"pivot": adminGroup.Pivot,
		})
	}

	Result := gin.H{
		"data": categoies,
	}

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