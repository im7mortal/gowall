package main

import (
	"github.com/im7mortal/gowall/handlers"
)



//func LoadTemplates(router *gin.Engine) {
func LoadTemplates() {

	defaultTmpl := "default.html"

	handlers.InitTemplate(defaultTmpl, "/", "layouts/" + defaultTmpl, "views/index.html")



}