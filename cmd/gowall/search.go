package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

func searchResult(c *gin.Context) {

	q, ok := c.GetQuery("q")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"users":          []int{},
			"accounts":       []int{},
			"administrators": []int{},
		})
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()

	users := make(chan []User, 1)
	accounts := make(chan []Account, 1)
	administrators := make(chan []Admin, 1)

	go func() {
		query := bson.M{
			"username": bson.RegEx{
				Pattern: `^.*?` + q + `.*$`,
				Options: "i",
			},
		}
		results := []User{}
		db.C(USERS).Find(query).Sort("username").Limit(10).All(&results)
		users <- results
	}()

	go func() {
		query := bson.M{
			"name.full": bson.RegEx{
				Pattern: `^.*?` + q + `.*$`,
				Options: "i",
			},
		}
		results := []Account{}
		db.C(ACCOUNTS).Find(query).Sort("name.full").Limit(10).All(&results)
		accounts <- results
	}()

	go func() {
		query := bson.M{
			"name.full": bson.RegEx{
				Pattern: `^.*?` + q + `.*$`,
				Options: "i",
			},
		}
		results := []Admin{}
		db.C(ADMINS).Find(query).Sort("name.full").Limit(10).All(&results)
		administrators <- results
	}()

	c.JSON(http.StatusOK, gin.H{
		"users":          <-users,
		"accounts":       <-accounts,
		"administrators": <-administrators,
	})
}

func getData(c *gin.Context, query *mgo.Query, results interface{}) (data gin.H) {
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
		"prev":    page - 1,
		"hasPrev": page - 1 != 0,
		"next":    page + 1,
		"hasNext": float64(count) / float64(count_) > 1,
		"total":   count,
	}

	end := count_
	if count_ > count {
		end = count
	}

	items := gin.H{
		"begin": (page - 1) * limit,
		"end":   end,
		"total": count,
	}

	filters := gin.H{
		"limit": limit,
		"page":  page,
		"sort":  sort,
	}
	return gin.H{
		"data":    results,
		"pages":   pages,
		"items":   items,
		"filters": filters,
	}
}
