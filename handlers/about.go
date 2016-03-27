package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func About(c *gin.Context) {
	c.Render(http.StatusOK, TemplateStorage[c.Request.URL.Path])
}
