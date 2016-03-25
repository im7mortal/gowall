package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func About(c *gin.Context) {
	//todo
	c.HTML(http.StatusOK, "default.html", gin.H{
		"title": "Main website",
		"pipeline": true,
	})
}
