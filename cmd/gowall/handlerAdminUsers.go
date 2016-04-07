package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
"gopkg.in/mgo.v2"
	"encoding/json"
)

func AdminUsersRender(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func UsersRender(c *gin.Context) {

	session, err := mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	collection := d.C(USERS)
	user := User{}
	collection.FindId(c.Param("id")).One(&user)
	userJSON, _ := json.Marshal(user)
	c.Set("Record", string(userJSON))
	render, _ := TemplateStorage["/admin/users/details/"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}
