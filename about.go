package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func About(c *gin.Context) {

	trigger = !trigger

	/*	if trigger {

			Router.SetHTMLTemplate(tm1)
			println(tm1.DefinedTemplates())
		} else {

			Router.SetHTMLTemplate(tm2)
			println(tm2.DefinedTemplates())
		}*/
	c.HTML(http.StatusOK, "default.html", gin.H{
		"title": "Main website",
		"pipeline": true,
	})
}
