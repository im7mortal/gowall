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
	InitTemplate(defaultTmpl, "/signup/social/", "views/signup/social.html")
	InitTemplate(defaultTmpl, "/login/", "views/login/index.html")
	InitTemplate(defaultTmpl, "/login/forgot/", "views/login/forgot/index.html")
	InitTemplate(defaultTmpl, "/login/reset/", "views/login/reset/index.html")
	InitTemplate(defaultTmpl, "404", "views/http/404.html")
	InitTemplate(defaultTmpl, "500", "views/http/500.html")

	InitTemplate(accountTmpl, "/account/", "views/account/index.html")
	InitTemplate(accountTmpl, "/account/verification/", "views/account/verification/index.html")
	InitTemplate(accountTmpl, "/account/settings/", "views/account/settings/index.html")

	InitTemplate(adminTmpl, "/admin/", "views/admin/index.html")
	InitTemplate(adminTmpl, "/admin/users/", "views/admin/users/index.html")
	InitTemplate(adminTmpl, "/admin/users/details/", "views/admin/users/details.html")
	InitTemplate(adminTmpl, "/admin/admin-groups/", "views/admin/admin-groups/index.html")
	InitTemplate(adminTmpl, "/admin/admin-groups/details/", "views/admin/admin-groups/details.html")
	InitTemplate(adminTmpl, "/admin/administrators/", "views/admin/administrators/index.html")
	InitTemplate(adminTmpl, "/admin/administrators/details/", "views/admin/administrators/details.html")
	InitTemplate(adminTmpl, "/admin/categories/", "views/admin/categories/index.html")
	InitTemplate(adminTmpl, "/admin/categories/details/", "views/admin/categories/details.html")
	InitTemplate(adminTmpl, "/admin/statuses/", "views/admin/statuses/index.html")
	InitTemplate(adminTmpl, "/admin/statuses/details/", "views/admin/statuses/details.html")

}