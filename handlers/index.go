package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.Render(http.StatusOK, TemplateStorage["/"])
}
