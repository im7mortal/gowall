package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func renderAdministrator(c *gin.Context) {
	db := getMongoDBInstance()
	defer db.Session.Close()

	CountAccountChan := make(chan int)
	CountUserChan := make(chan int)
	CountAdminChan := make(chan int)
	CountAdminGroupChan := make(chan int)
	CountCategoryChan := make(chan int)
	CountStatusChan := make(chan int)

	go getCount(db.C(ACCOUNTS), CountAccountChan, bson.M{})
	go getCount(db.C(USERS), CountUserChan, bson.M{})
	go getCount(db.C(ADMINS), CountAdminChan, bson.M{})
	go getCount(db.C(ADMINGROUPS), CountAdminGroupChan, bson.M{})
	go getCount(db.C(CATEGORIES), CountCategoryChan, bson.M{})
	go getCount(db.C(STATUSES), CountStatusChan, bson.M{})

	c.Set("CountAccount", <-CountAccountChan)
	c.Set("CountUser", <-CountUserChan)
	c.Set("CountAdmin", <-CountAdminChan)
	c.Set("CountAdminGroup", <-CountAdminGroupChan)
	c.Set("CountCategory", <-CountCategoryChan)
	c.Set("CountStatus", <-CountStatusChan)

	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}
