package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"html/template"
	"strings"

	"gopkg.in/mgo.v2"
	"net/url"
)

func AdminUsersRender(c *gin.Context) {
	query := bson.M{}

	username, ok := c.GetQuery("username")
	if ok && len(username) != 0 {
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

	type _user struct {
		ID bson.ObjectId `bson:"_id" json:"_id"`
		Username string `bson:"username" json:"username"`
		IsActive string `bson:"isActive" json:"isActive"`
		Email string `bson:"email" json:"email"`
	}
	var results []_user

	db := getMongoDBInstance()
	defer db.Session.Close()
	// TODO keys
	collection := db.C(USERS)

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

func getData (c *gin.Context, query *mgo.Query, results interface{}) (data gin.H) {
	limitS := c.DefaultQuery("limit", "20")
	limit_, _ := strconv.ParseInt(limitS, 0, 0)
	limit := int(limit_)
	if limit > 100 {
		limit = 100
	}

	pageS := c.DefaultQuery("page", "0")
	page_, _ := strconv.ParseInt(pageS, 0, 0)
	page := int(page_)
	sort := c.DefaultQuery("sort", "_id")

	count, _ := query.Count()
	query.Skip(page * limit).Sort(sort).Limit(limit).All(results)

	page += 1
	count_ := page * limit
	pages := gin.H{
		"current": page,
		"prev": page - 1,
		"hasPrev": page - 1 != 0,
		"next": page + 1,
		"hasNext": float64(count) / float64(count_) > 1,
		"total": count,
	}

	end := count_
	if count_ > count {
		end = count
	}

	items := gin.H{
		"begin": (page - 1) * limit,
		"end": end,
		"total": count,
	}

	return gin.H{
		"data": results,
		"pages": pages,
		"items": items,
	}
}


func CreateUser(c *gin.Context) {
	response := Response{} // todo sync.Pool
	defer response.Recover(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
		return
	}
	// clean errors from client
	response.CleanErrors()

	// validate
	response.ValidateUsername()

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"username": response.Username}).One(&user)
	if err != nil {
		println(err.Error())
	}

	// duplicateUsernameCheck
	if len(user.Username) != 0 {
		if user.Username == response.Username {
			response.Errors = append(response.Errors, "That username is already taken.")
		}
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	// createUser

	user.ID = bson.NewObjectId()
	user.Username = response.Username
	user.Search = []string{response.Username}

	err = collection.Insert(user)
	if err != nil {
		panic(err)
		return
	}
}

func UsersRender(c *gin.Context) {

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	userJSON, _ := json.Marshal(gin.H{
		"_id": user.ID.Hex(),
		"username": user.Username,
		"email": user.Email,
		"isActive": user.IsActive,
	})
	c.Set("Record", template.JS(url.QueryEscape(string(userJSON))))
	render, _ := TemplateStorage["/admin/users/details/"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

