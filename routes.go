package main

import "github.com/gin-gonic/gin"

func BindRoutes(router *gin.Engine) {

	//front end
	//router.GET("/", Smock)
	router.GET("/about/", About)
	router.GET("/contact/", About)
	router.POST("/contact/", About)

	//sign up
	router.GET("/signup/", About)
	router.POST("/signup/", About)

	//social sign up
	router.POST("/signup/social/", About)
	router.GET("/signup/twitter/", About)
	router.GET("/signup/twitter/callback/", About)
	router.GET("/signup/github/", About)
	router.GET("/signup/github/callback/", About)
	router.GET("/signup/facebook/", About)
	router.GET("/signup/facebook/callback/", About)
	router.GET("/signup/google/", About)
	router.GET("/signup/google/callback/", About)
	router.GET("/signup/tumblr/", About)
	router.GET("/signup/tumblr/callback/", About)

	//login/out
	router.GET("/login/", About)
	router.POST("/login/", About)
	router.GET("/login/forgot/", About)
	router.POST("/login/forgot/", About)
	router.GET("/login/reset/", About)
	router.GET("/login/reset/:email/:token/", About)
	router.PUT("/login/reset/:email/:token/", About)
	router.GET("/logout/", About)

	//social login
	router.GET("/login/twitter/", About)
	router.GET("/login/twitter/callback/", About)
	router.GET("/login/github/", About)
	router.GET("/login/github/callback/", About)
	router.GET("/login/facebook/", About)
	router.GET("/login/facebook/callback/", About)
	router.GET("/login/google/", About)
	router.GET("/login/google/callback/", About)
	router.GET("/login/tumblr/", About)
	router.GET("/login/tumblr/callback/", About)

	//admin
	//app.all('/admin*', ensureAuthenticated); TODO necessary midleware
	//app.all('/admin*', ensureAdmin);
	router.GET("/admin/", About)

	//admin > users
	router.GET("/admin/users/", About)
	router.POST("/admin/users/", About)
	router.GET("/admin/users/:id/", About)
	router.PUT("/admin/users/:id/", About)
	router.PUT("/admin/users/:id/password/", About)
	router.PUT("/admin/users/:id/role-admin/", About)
	router.DELETE("/admin/users/:id/role-admin/", About)
	router.PUT("/admin/users/:id/role-account/", About)
	router.DELETE("/admin/users/:id/role-account/", About)
	router.DELETE("/admin/users/:id/", About)

	//admin > administrators
	router.GET("/admin/administrators/", About)
	router.POST("/admin/administrators/", About)
	router.GET("/admin/administrators/:id/", About)
	router.PUT("/admin/administrators/:id/", About)
	router.PUT("/admin/administrators/:id/permissions/", About)
	router.PUT("/admin/administrators/:id/groups/", About)
	router.PUT("/admin/administrators/:id/user/", About)
	router.DELETE("/admin/administrators/:id/user/", About)
	router.DELETE("/admin/administrators/:id/", About)

	//admin > admin groups
	router.GET("/admin/admin-groups/", About)
	router.POST("/admin/admin-groups/", About)
	router.GET("/admin/admin-groups/:id/", About)
	router.PUT("/admin/admin-groups/:id/", About)
	router.PUT("/admin/admin-groups/:id/permissions/", About)
	router.DELETE("/admin/admin-groups/:id/", About)

	//admin > accounts
	router.GET("/admin/accounts/", About)
	router.POST("/admin/accounts/", About)
	router.GET("/admin/accounts/:id/", About)
	router.PUT("/admin/accounts/:id/", About)
	router.PUT("/admin/accounts/:id/user/", About)
	router.DELETE("/admin/accounts/:id/user/", About)
	router.POST("/admin/accounts/:id/notes/", About)
	router.POST("/admin/accounts/:id/status/", About)
	router.DELETE("/admin/accounts/:id/", About)


	//admin > statuses
	router.GET("/admin/statuses/", About)
	router.POST("/admin/statuses/", About)
	router.GET("/admin/statuses/:id/", About)
	router.PUT("/admin/statuses/:id/", About)
	router.DELETE("/admin/statuses/:id/", About)


	//admin > categories
	router.GET("/admin/categories/", About)
	router.POST("/admin/categories/", About)
	router.GET("/admin/categories/:id/", About)
	router.PUT("/admin/categories/:id/", About)
	router.DELETE("/admin/categories/:id/", About)

	//admin > search
	router.GET("/admin/search/", About)

	//account
	//app.all('/account*', ensureAuthenticated); //TODO necessary midleware
	//app.all('/account*', ensureAccount);
	router.GET("/account/", About)

	//account > verification
	router.GET("/account/verification/", About)
	router.POST("/account/verification/", About)
	router.GET("/account/verification/:token/", About)

	//account > settings
	router.GET("/account/settings/", About)
	router.PUT("/account/settings/", About)
	router.PUT("/account/settings/identity/", About)
	router.PUT("/account/settings/password/", About)

	//account > settings > social
	router.GET("/account/settings/twitter/", About)
	router.GET("/account/settings/twitter/callback/", About)
	router.GET("/account/settings/twitter/disconnect/", About)
	router.GET("/account/settings/github/", About)
	router.GET("/account/settings/github/callback/", About)
	router.GET("/account/settings/github/disconnect/", About)
	router.GET("/account/settings/facebook/", About)
	router.GET("/account/settings/facebook/callback/", About)
	router.GET("/account/settings/facebook/disconnect/", About)
	router.GET("/account/settings/google/", About)
	router.GET("/account/settings/google/callback/", About)
	router.GET("/account/settings/google/disconnect/", About)
	router.GET("/account/settings/tumblr/", About)
	router.GET("/account/settings/tumblr/callback/", About)
	router.GET("/account/settings/tumblr/disconnect/", About)

	//route not found
	router.NoRoute(About)
}