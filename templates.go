package main

import (
	"github.com/im7mortal/gowall/handlers"
)



//func LoadTemplates(router *gin.Engine) {
func LoadTemplates() {

	defaultTmpl := "default.html"
	accountTmpl := "account.html"

	handlers.InitTemplate(defaultTmpl, "/", "views/index.html")
	handlers.InitTemplate(defaultTmpl, "/about/", "views/about/index.html")
	handlers.InitTemplate(defaultTmpl, "/contact/", "views/contact/index.html")
	handlers.InitTemplate(defaultTmpl, "/signup/", "views/signup/index.html")
	handlers.InitTemplate(defaultTmpl, "/login/", "views/login/index.html")
	handlers.InitTemplate(defaultTmpl, "/login/forgot/", "views/login/forgot/index.html")
	handlers.InitTemplate(defaultTmpl, "404", "views/http/404.html")

	handlers.InitTemplate(accountTmpl, "/account/", "views/account/index.html")
	handlers.InitTemplate(accountTmpl, "/account/verification/", "views/account/verification/index.html")
	handlers.InitTemplate(accountTmpl, "/account/settings/", "views/account/settings/index.html")



}