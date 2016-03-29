package main

import (
	"github.com/im7mortal/gowall/handlers"
)



//func LoadTemplates(router *gin.Engine) {
func LoadTemplates() {

	defaultTmpl := "default.html"

	handlers.InitTemplate(defaultTmpl, "/", "views/index.html")
	handlers.InitTemplate(defaultTmpl, "/about/", "views/about/index.html")
	handlers.InitTemplate(defaultTmpl, "/contact/", "views/contact/index.html")
	handlers.InitTemplate(defaultTmpl, "/signup/", "views/signup/index.html")



}