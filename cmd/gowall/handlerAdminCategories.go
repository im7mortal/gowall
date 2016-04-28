package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"strings"
	"regexp"
	"gopkg.in/mgo.v2"
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

	var results []Category

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

func CreateCategory(c *gin.Context) {
	response := Response{} // todo sync.Pool
	defer response.Recover(c)

	admin, ok := getAdmin(c)
	if !ok {// todo extra
		panic("not authorised")
	}

	// validate
	ok = admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not create categories")
		response.Fail(c)
		return
	}

	category := Category{}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&category)
	if err != nil {
		panic(err)
		return
	}
	// clean errors from client
	response.CleanErrors()

	if len(category.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if len(category.Pivot) == 0 {
		response.Errors = append(response.Errors, "A pivot is required")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	//duplicateCategoryCheck
	_id := slugify(category.Pivot + " " + category.Name)
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(CATEGORIES)
	category_ := Category{}
	err = collection.Find(bson.M{"_id": _id}).One(&category_)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That category+pivot is already taken.")
		response.Fail(c)
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
	response.Success = true
	c.JSON(http.StatusOK, response)
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

// TODO add g greedy
var r1, _ = regexp.Compile(`[^\w ]+`)
var r2, _ = regexp.Compile(` +`)

func slugify(str string) string {
	str = strings.ToLower(str)
	str_ := []byte(str)
	str_ = r1.ReplaceAll(str_, []byte(""))
	return string(r2.ReplaceAll(str_, []byte("-")))
}