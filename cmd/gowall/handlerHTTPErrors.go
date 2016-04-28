package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func Status404Render(c *gin.Context) {
	render, _ := TemplateStorage["404"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
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
				render, _ := TemplateStorage["500"]
				render.Data = gin.H{
					"Stack": fmt.Sprintf("%v\n", rec),
				}
				c.Render(http.StatusInternalServerError, render)
			}
		}
	}(c)
	c.Next()
}
