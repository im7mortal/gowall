package main

import "github.com/gin-gonic/gin"

func BindRoutes(router *gin.Engine) {

	//front end
	router.GET("/", Smock)
	router.GET("/about/", Smock)
	router.GET("/contact/", Smock)
	router.POST("/contact/", Smock)

	//sign up
	router.GET("/signup/", Smock)
	router.POST("/signup/", Smock)

	//social sign up
	router.POST("/signup/social/", Smock)
	router.GET("/signup/twitter/", Smock)
	router.GET("/signup/twitter/callback/", Smock)
	router.GET("/signup/github/", Smock)
	router.GET("/signup/github/callback/", Smock)
	router.GET("/signup/facebook/", Smock)
	router.GET("/signup/facebook/callback/", Smock)
	router.GET("/signup/google/", Smock)
	router.GET("/signup/google/callback/", Smock)
	router.GET("/signup/tumblr/", Smock)
	router.GET("/signup/tumblr/callback/", Smock)

	//login/out
	router.GET("/login/", Smock)
	router.POST("/login/", Smock)
	router.GET("/login/forgot/", Smock)
	router.POST("/login/forgot/", Smock)
	router.GET("/login/reset/", Smock)
	router.GET("/login/reset/:email/:token/", Smock)
	router.PUT("/login/reset/:email/:token/", Smock)
	router.GET("/logout/", Smock)

	//social login
	router.GET("/login/twitter/", Smock)
	router.GET("/login/twitter/callback/", Smock)
	router.GET("/login/github/", Smock)
	router.GET("/login/github/callback/", Smock)
	router.GET("/login/facebook/", Smock)
	router.GET("/login/facebook/callback/", Smock)
	router.GET("/login/google/", Smock)
	router.GET("/login/google/callback/", Smock)
	router.GET("/login/tumblr/", Smock)
	router.GET("/login/tumblr/callback/", Smock)

	//admin
	//app.all('/admin*', ensureAuthenticated); TODO necessary midleware
	//app.all('/admin*', ensureAdmin);
	router.GET("/admin/", Smock)

	//admin > users
	router.GET("/admin/users/", Smock)
	router.POST("/admin/users/", Smock)
	router.GET("/admin/users/:id/", Smock)
	router.PUT("/admin/users/:id/", Smock)
	router.PUT("/admin/users/:id/password/", Smock)
	router.PUT("/admin/users/:id/role-admin/", Smock)
	router.DELETE("/admin/users/:id/role-admin/", Smock)
	router.PUT("/admin/users/:id/role-account/", Smock)
	router.DELETE("/admin/users/:id/role-account/", Smock)
	router.DELETE("/admin/users/:id/", Smock)

	//admin > administrators
	router.GET("/admin/administrators/", Smock)
	router.POST("/admin/administrators/", Smock)
	router.GET("/admin/administrators/:id/", Smock)
	router.PUT("/admin/administrators/:id/", Smock)
	router.PUT("/admin/administrators/:id/permissions/", Smock)
	router.PUT("/admin/administrators/:id/groups/", Smock)
	router.PUT("/admin/administrators/:id/user/", Smock)
	router.DELETE("/admin/administrators/:id/user/", Smock)
	router.DELETE("/admin/administrators/:id/", Smock)

	//admin > admin groups
	router.GET("/admin/admin-groups/", Smock)
	router.POST("/admin/admin-groups/", Smock)
	router.GET("/admin/admin-groups/:id/", Smock)
	router.PUT("/admin/admin-groups/:id/", Smock)
	router.PUT("/admin/admin-groups/:id/permissions/", Smock)
	router.DELETE("/admin/admin-groups/:id/", Smock)

	//admin > accounts
	router.GET("/admin/accounts/", Smock)
	router.POST("/admin/accounts/", Smock)
	router.GET("/admin/accounts/:id/", Smock)
	router.PUT("/admin/accounts/:id/", Smock)
	router.PUT("/admin/accounts/:id/user/", Smock)
	router.DELETE("/admin/accounts/:id/user/", Smock)
	router.POST("/admin/accounts/:id/notes/", Smock)
	router.POST("/admin/accounts/:id/status/", Smock)
	router.DELETE("/admin/accounts/:id/", Smock)


	//admin > statuses
	router.GET("/admin/statuses/", Smock)
	router.POST("/admin/statuses/", Smock)
	router.GET("/admin/statuses/:id/", Smock)
	router.PUT("/admin/statuses/:id/", Smock)
	router.DELETE("/admin/statuses/:id/", Smock)


	//admin > categories
	router.GET("/admin/categories/", Smock)
	router.POST("/admin/categories/", Smock)
	router.GET("/admin/categories/:id/", Smock)
	router.PUT("/admin/categories/:id/", Smock)
	router.DELETE("/admin/categories/:id/", Smock)

	//admin > search
	router.GET("/admin/search/", Smock)

	//account
	//app.all('/account*', ensureAuthenticated); //TODO necessary midleware
	//app.all('/account*', ensureAccount);
	router.GET("/account/", Smock)

	//account > verification
	router.GET("/account/verification/", Smock)
	router.POST("/account/verification/", Smock)
	router.GET("/account/verification/:token/", Smock)

	//account > settings
	router.GET("/account/settings/", Smock)
	router.PUT("/account/settings/", Smock)
	router.PUT("/account/settings/identity/", Smock)
	router.PUT("/account/settings/password/", Smock)

	//account > settings > social
	router.GET("/account/settings/twitter/", Smock)
	router.GET("/account/settings/twitter/callback/", Smock)
	router.GET("/account/settings/twitter/disconnect/", Smock)
	router.GET("/account/settings/github/", Smock)
	router.GET("/account/settings/github/callback/", Smock)
	router.GET("/account/settings/github/disconnect/", Smock)
	router.GET("/account/settings/facebook/", Smock)
	router.GET("/account/settings/facebook/callback/", Smock)
	router.GET("/account/settings/facebook/disconnect/", Smock)
	router.GET("/account/settings/google/", Smock)
	router.GET("/account/settings/google/callback/", Smock)
	router.GET("/account/settings/google/disconnect/", Smock)
	router.GET("/account/settings/tumblr/", Smock)
	router.GET("/account/settings/tumblr/callback/", Smock)
	router.GET("/account/settings/tumblr/disconnect/", Smock)

	//route not found
	router.NoRoute(Smock)
}