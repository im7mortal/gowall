package main

import (
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
)

func initTemplates(router *gin.Engine) (r multitemplate.Render) {

	r = multitemplate.New()

	default_ := tmpl{"layouts/default.html", r}
	account := tmpl{"layouts/account.html", r}
	admin := tmpl{"layouts/admin.html", r}

	default_.addFromFiles("/", "views/index.html")
	default_.addFromFiles("/about/", "views/about/index.html")
	default_.addFromFiles("/contact/", "views/contact/index.html")
	default_.addFromFiles("/signup/", "views/signup/index.html")
	default_.addFromFiles("/signup/social/", "views/signup/social.html")
	default_.addFromFiles("/login/", "views/login/index.html")
	default_.addFromFiles("/login/forgot/", "views/login/forgot/index.html")
	default_.addFromFiles("/login/reset/", "views/login/reset/index.html")
	default_.addFromFiles("404", "views/http/404.html")
	default_.addFromFiles("500", "views/http/500.html")

	account.addFromFiles("/account/", "views/account/index.html")
	account.addFromFiles("/account/verification/", "views/account/verification/index.html")
	account.addFromFiles("/account/settings/", "views/account/settings/index.html")

	admin.addFromFiles("/admin/", "views/admin/index.html")
	admin.addFromFiles("/admin/users/", "views/admin/users/index.html")
	admin.addFromFiles("/admin/users/details/", "views/admin/users/details.html")
	admin.addFromFiles("/admin/admin-groups/", "views/admin/admin-groups/index.html")
	admin.addFromFiles("/admin/admin-groups/details/", "views/admin/admin-groups/details.html")
	admin.addFromFiles("/admin/administrators/", "views/admin/administrators/index.html")
	admin.addFromFiles("/admin/administrators/details/", "views/admin/administrators/details.html")
	admin.addFromFiles("/admin/categories/", "views/admin/categories/index.html")
	admin.addFromFiles("/admin/categories/details/", "views/admin/categories/details.html")
	admin.addFromFiles("/admin/statuses/", "views/admin/statuses/index.html")
	admin.addFromFiles("/admin/statuses/details/", "views/admin/statuses/details.html")
	return
}


type tmpl struct {
	rootTmpl string
	r multitemplate.Render
}

func (t *tmpl) addFromFiles (name string, files ...string) {
	files = append([]string{t.rootTmpl}, files...)
	t.r.AddFromFiles(name, files...)
}
