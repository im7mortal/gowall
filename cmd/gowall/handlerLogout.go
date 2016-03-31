package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	removedCookie := http.Cookie{
		Name:  "session",
		MaxAge:   -1,
	}
	http.SetCookie(c.Writer, &removedCookie)
	c.Redirect(http.StatusFound, "/home/")
}