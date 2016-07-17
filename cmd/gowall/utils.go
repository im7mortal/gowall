package main

import (
	"net/url"
	"strings"
	"regexp"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"strconv"
)

func getEscapedString(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

var r1, _ = regexp.Compile(`[^\w ]+`)
var r2, _ = regexp.Compile(` +`)

/**
	preparing id
 */

func slugify (str string) string {
	str = strings.ToLower(str)
	str = r1.ReplaceAllString(str, "")
	str = r2.ReplaceAllString(str, "-")
	return str
}

func slugifyName(str string) string {
	str = strings.TrimSpace(str)
	return r2.ReplaceAllString(str, " ")
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

	filters := gin.H{
		"limit": limit,
		"page": page,
		"sort": sort,
	}
	return gin.H{
		"data": results,
		"pages": pages,
		"items": items,
		"filters": filters,
	}
}



func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}