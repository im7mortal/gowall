package main

import (
	"github.com/gin-gonic/gin"
	"github.com/im7mortal/gowall/handlers"
)

func BindRoutes(router *gin.Engine) {

	//front end
	router.GET("/", handlers.Index)
	router.GET("/about/", handlers.About)
	router.GET("/contact/", handlers.ContactRender)
	router.POST("/contact/", handlers.ContactSend)

	//sign up
	router.GET("/signup/", handlers.SignupRender)
	router.POST("/signup/", handlers.Signup)

	//social sign up
	router.POST("/signup/social/", handlers.Index)
	router.GET("/signup/twitter/", handlers.Index)
	router.GET("/signup/twitter/callback/", handlers.Index)
	router.GET("/signup/github/", handlers.Index)
	router.GET("/signup/github/callback/", handlers.Index)
	router.GET("/signup/facebook/", handlers.Index)
	router.GET("/signup/facebook/callback/", handlers.Index)
	router.GET("/signup/google/", handlers.Index)
	router.GET("/signup/google/callback/", handlers.Index)
	router.GET("/signup/tumblr/", handlers.Index)
	router.GET("/signup/tumblr/callback/", handlers.Index)

	//login/out
	router.GET("/login/", handlers.LoginRender)
	router.POST("/login/", handlers.Login)
	router.GET("/login/forgot/", handlers.ForgotRender)
	router.POST("/login/forgot/", handlers.Index)
	router.GET("/login/reset/", handlers.Index)
	router.GET("/login/reset/:email/:token/", handlers.Index)
	router.PUT("/login/reset/:email/:token/", handlers.Index)
	router.GET("/logout/", handlers.Logout)

	//social login
	router.GET("/login/twitter/", handlers.Index)
	router.GET("/login/twitter/callback/", handlers.Index)
	router.GET("/login/github/", handlers.Index)
	router.GET("/login/github/callback/", handlers.Index)
	router.GET("/login/facebook/", handlers.Index)
	router.GET("/login/facebook/callback/", handlers.Index)
	router.GET("/login/google/", handlers.Index)
	router.GET("/login/google/callback/", handlers.Index)
	router.GET("/login/tumblr/", handlers.Index)
	router.GET("/login/tumblr/callback/", handlers.Index)

	//admin
	//app.all('/admin*', ensureAuthenticated); TODO necessary midleware
	//app.all('/admin*', ensureAdmin);
	router.GET("/admin/", handlers.Index)

	//admin > users
	router.GET("/admin/users/", handlers.Index)
	router.POST("/admin/users/", handlers.Index)
	router.GET("/admin/users/:id/", handlers.Index)
	router.PUT("/admin/users/:id/", handlers.Index)
	router.PUT("/admin/users/:id/password/", handlers.Index)
	router.PUT("/admin/users/:id/role-admin/", handlers.Index)
	router.DELETE("/admin/users/:id/role-admin/", handlers.Index)
	router.PUT("/admin/users/:id/role-account/", handlers.Index)
	router.DELETE("/admin/users/:id/role-account/", handlers.Index)
	router.DELETE("/admin/users/:id/", handlers.Index)

	//admin > administrators
	router.GET("/admin/administrators/", handlers.Index)
	router.POST("/admin/administrators/", handlers.Index)
	router.GET("/admin/administrators/:id/", handlers.Index)
	router.PUT("/admin/administrators/:id/", handlers.Index)
	router.PUT("/admin/administrators/:id/permissions/", handlers.Index)
	router.PUT("/admin/administrators/:id/groups/", handlers.Index)
	router.PUT("/admin/administrators/:id/user/", handlers.Index)
	router.DELETE("/admin/administrators/:id/user/", handlers.Index)
	router.DELETE("/admin/administrators/:id/", handlers.Index)

	//admin > admin groups
	router.GET("/admin/admin-groups/", handlers.Index)
	router.POST("/admin/admin-groups/", handlers.Index)
	router.GET("/admin/admin-groups/:id/", handlers.Index)
	router.PUT("/admin/admin-groups/:id/", handlers.Index)
	router.PUT("/admin/admin-groups/:id/permissions/", handlers.Index)
	router.DELETE("/admin/admin-groups/:id/", handlers.Index)

	//admin > accounts
	router.GET("/admin/accounts/", handlers.Index)
	router.POST("/admin/accounts/", handlers.Index)
	router.GET("/admin/accounts/:id/", handlers.Index)
	router.PUT("/admin/accounts/:id/", handlers.Index)
	router.PUT("/admin/accounts/:id/user/", handlers.Index)
	router.DELETE("/admin/accounts/:id/user/", handlers.Index)
	router.POST("/admin/accounts/:id/notes/", handlers.Index)
	router.POST("/admin/accounts/:id/status/", handlers.Index)
	router.DELETE("/admin/accounts/:id/", handlers.Index)


	//admin > statuses
	router.GET("/admin/statuses/", handlers.Index)
	router.POST("/admin/statuses/", handlers.Index)
	router.GET("/admin/statuses/:id/", handlers.Index)
	router.PUT("/admin/statuses/:id/", handlers.Index)
	router.DELETE("/admin/statuses/:id/", handlers.Index)


	//admin > categories
	router.GET("/admin/categories/", handlers.Index)
	router.POST("/admin/categories/", handlers.Index)
	router.GET("/admin/categories/:id/", handlers.Index)
	router.PUT("/admin/categories/:id/", handlers.Index)
	router.DELETE("/admin/categories/:id/", handlers.Index)

	//admin > search
	router.GET("/admin/search/", handlers.Index)

	//account
	//app.all('/account*', ensureAuthenticated); //TODO necessary midleware
	//app.all('/account*', ensureAccount);
	router.GET("/account/", handlers.Account)

	//account > verification
	router.GET("/account/verification/", handlers.AccountVerification)
	router.POST("/account/verification/", handlers.Index)
	router.GET("/account/verification/:token/", handlers.Index)

	//account > settings
	router.GET("/account/settings/", handlers.AccountSettingsRender)
	router.PUT("/account/settings/", handlers.Index)
	router.PUT("/account/settings/identity/", handlers.Index)
	router.PUT("/account/settings/password/", handlers.Index)

	//account > settings > social
	router.GET("/account/settings/twitter/", handlers.Index)
	router.GET("/account/settings/twitter/callback/", handlers.Index)
	router.GET("/account/settings/twitter/disconnect/", handlers.Index)
	router.GET("/account/settings/github/", handlers.Index)
	router.GET("/account/settings/github/callback/", handlers.Index)
	router.GET("/account/settings/github/disconnect/", handlers.Index)
	router.GET("/account/settings/facebook/", handlers.Index)
	router.GET("/account/settings/facebook/callback/", handlers.Index)
	router.GET("/account/settings/facebook/disconnect/", handlers.Index)
	router.GET("/account/settings/google/", handlers.Index)
	router.GET("/account/settings/google/callback/", handlers.Index)
	router.GET("/account/settings/google/disconnect/", handlers.Index)
	router.GET("/account/settings/tumblr/", handlers.Index)
	router.GET("/account/settings/tumblr/callback/", handlers.Index)
	router.GET("/account/settings/tumblr/disconnect/", handlers.Index)

	//route not found
	router.NoRoute(handlers.Index)
}