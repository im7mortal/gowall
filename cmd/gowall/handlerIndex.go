package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	println(c.Request.Host)
	render, _ := TemplateStorage["/"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}
