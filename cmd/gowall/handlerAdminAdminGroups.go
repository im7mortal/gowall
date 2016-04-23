package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"html/template"
)

func AdminAdminGroupsRender(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.M{
			"$regex": "/^.*?" + name + ".*$/i",
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

	var results []AdminGroup

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINGROUPS)
	collection.Find(query).Skip(int(limit * page)).Sort(sort).Limit(int(limit)).All(&results)

	users := []gin.H{}

	for _, adminGroup := range results {
		users = append(users, gin.H{
			"_id": adminGroup.ID.Hex(),
			"name": adminGroup.Name,
		})
	}

	Result := gin.H{
		"data": users,
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

func AdminGroupRender(c *gin.Context) {

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
