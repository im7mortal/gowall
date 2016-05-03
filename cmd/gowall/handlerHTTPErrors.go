package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func Status404Render(c *gin.Context) {
	c.HTML(http.StatusOK, "404", c.Keys)
}

func checkRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			if XHR(c) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Something went wrong.",
					"details": rec,
				})
			} else {
				c.Set("Stack", fmt.Sprintf("%v\n", rec))
				c.HTML(http.StatusOK, "500", c.Keys)
			}
		}
	}(c)
	c.Next()
}
