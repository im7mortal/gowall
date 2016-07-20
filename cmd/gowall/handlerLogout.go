package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("public")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
