package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, TemplateStorage["/"])
}
