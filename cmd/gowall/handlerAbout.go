package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func About(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}
