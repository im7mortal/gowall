package main

//func LoadTemplates(router *gin.Engine) {
func LoadTemplates() {

	defaultTmpl := "default.html"
	accountTmpl := "account.html"
	adminTmpl := "admin.html"

	InitTemplate(defaultTmpl, "/", "views/index.html")
	InitTemplate(defaultTmpl, "/about/", "views/about/index.html")
	InitTemplate(defaultTmpl, "/contact/", "views/contact/index.html")
	InitTemplate(defaultTmpl, "/signup/", "views/signup/index.html")
	InitTemplate(defaultTmpl, "/login/", "views/login/index.html")
	InitTemplate(defaultTmpl, "/login/forgot/", "views/login/forgot/index.html")
	InitTemplate(defaultTmpl, "404", "views/http/404.html")

	InitTemplate(accountTmpl, "/account/", "views/account/index.html")
	InitTemplate(accountTmpl, "/account/verification/", "views/account/verification/index.html")
	InitTemplate(accountTmpl, "/account/settings/", "views/account/settings/index.html")

	InitTemplate(adminTmpl, "/admin/", "views/admin/index.html")



}