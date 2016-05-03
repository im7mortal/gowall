package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func About(c *gin.Context) {
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}
