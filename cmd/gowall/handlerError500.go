package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Status404Render(c *gin.Context) {
	render, _ := TemplateStorage["404"]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}
