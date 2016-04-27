package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"html/template"
	"strings"
)

func AdminUsersRender(c *gin.Context) {
	query := bson.M{}

	username, ok := c.GetQuery("username")
	if ok && len(username) != 0 {
		/*query["username"] = bson.M{
			"$regex": "/^.*?" + username + ".*$/i",
		}*/
		query["username"] = bson.RegEx{
			Pattern: `^.*?` + username + `.*$`,
			Options: "i",
		}
	}

	isActive, ok := c.GetQuery("isActive")
	if ok && len(isActive) != 0 {
		query["isActive"] = isActive
	}


	roles, ok := c.GetQuery("roles")
	if ok && len(roles) != 0 {
		// roles.admin or roles.account
		query["roles." + roles] = bson.M{
			"$exists": true,
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

	var results []User

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)

	collection.Find(query).Skip(int(limit * page)).Sort(sort).Limit(int(limit)).All(&results)

	users := []gin.H{}

	for _, user := range results {
		users = append(users, gin.H{
			"_id": user.ID.Hex(),
			"username": user.Username,
			"isActive": user.IsActive,
			"email": user.Email,
		})
	}

	Result := gin.H{
		"data": users,
		"pages": gin.H{
			"current": page + 1,
			"prev": 0,
			"hasPrev": false,
			"next": 0,
			"hasNext": false,
			"total": 0,
		},
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

func UsersRender(c *gin.Context) {

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

func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
