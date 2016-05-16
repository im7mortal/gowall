package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"gopkg.in/mgo.v2"
	"strings"
)

type responseAdmin struct {
	Response
	Admin
}

func renderAdministrators(c *gin.Context) {
	query := bson.M{}

	name, ok := c.GetQuery("name")
	if ok && len(name) != 0 {
		query["name"] = bson.RegEx{
			Pattern: `^.*?` + name + `.*$`,
			Options: "i",
		}
	}

	var results []Admin

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	Result := getData(c, collection.Find(query), &results)

	Results, err := json.Marshal(Result)
	if err != nil {
		panic(err)
	}
	if XHR(c) {
		handleXHR(c, Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func createAdministrator(c *gin.Context) {
	response := responseAdmin{}

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not create administrators")
		response.Fail(c)
		return
	}

	response.Admin.DecodeRequest(c)
	if len(response.Name.Full) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	// handleName
	response.Name.Full = slugifyName(response.Name.Full)

	// duplicateAdministrator
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	err := collection.Find(bson.M{"name.full": response.Name.Full}).One(nil)
	// we expect err == mgo.ErrNotFound for success
	if err == nil {
		response.Errors = append(response.Errors, "That administrator already exists.")
		response.Fail(c)
		return
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	// handleName
	name := strings.Split(response.Name.Full, " ")
	response.Name.First = name[0]
	if len(name) == 2 {
		response.Name.Last = name[1]
	}
	if len(name) == 3 {
		response.Name.Middle = name[2]
	}

	// createAdministrator
	response.Admin.ID = bson.NewObjectId()
	err = collection.Insert(response.Admin) // todo I think mgo's behavior isn't expected
	if err != nil {
		panic(err)
		return
	}
	response.Success = true
	c.JSON(http.StatusOK, gin.H{"record": response, "success": true})
}

func readAdministrator(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	admin := Admin{}
	err := collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&admin)
	if err != nil {
		if err == mgo.ErrNotFound {
			Status404Render(c)
			return
		}
		panic(err)
	}
	json, err := json.Marshal(admin)
	if err != nil {
		panic(err)
	}
	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(getEscapedString(string(json))))
	c.HTML(http.StatusOK, "/admin/administrators/details/", c.Keys)
}


func updateAdministrator(c *gin.Context) {
	response := responseAdmin{}
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

/*	if len(response.Name) == 0 {
		response.Errors = append(response.Errors, "A name is required")
	}*/

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)

	// patchAdminGroup
	err = collection.UpdateId(bson.ObjectIdHex(c.Param("id")), response.Admin)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}


	response.Success = true
	c.JSON(http.StatusOK, response)
}

func updateAdministratorPermissions(c *gin.Context) {
	response := responseAdmin{}
	defer response.Recover(c)

	admin := getAdmin(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not change the permissions of admin groups.")
		response.Fail(c)
		return
	}

	response.Admin.DecodeRequest(c)

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
	collection := db.C(ADMINS)

	err := collection.UpdateId(c.Param("id"), response.Admin)
	if err != nil {
		panic(err)
	}

	response.Finish(c)
}

func deleteAdministrator(c *gin.Context) {
	response := Response{} // todo sync.Pool

	// validate
	if ok := getAdmin(c).IsMemberOf("root"); !ok {
		response.Errors = append(response.Errors, "You may not delete administrators.")
		response.Fail(c)
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(ADMINS)
	err := collection.RemoveId(bson.ObjectIdHex(c.Param("id")))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	response.Finish(c)
}
