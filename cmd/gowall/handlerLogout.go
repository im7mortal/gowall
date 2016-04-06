package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/contrib/sessions"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("public")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
