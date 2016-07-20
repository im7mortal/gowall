package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "/", c.Keys)
}
