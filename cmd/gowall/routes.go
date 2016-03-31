package main

import (
	"github.com/gin-gonic/gin"
)

func BindRoutes(router *gin.Engine) {

	//front end
	router.GET("/", Index)
	router.GET("/about/", About)
	router.GET("/contact/", ContactRender)
	router.POST("/contact/", ContactSend)

	//sign up
	router.GET("/signup/", SignupRender)
	router.POST("/signup/", Signup)

	//social sign up
	router.POST("/signup/social/", Index)
	router.GET("/signup/twitter/", Index)
	router.GET("/signup/twitter/callback/", Index)
	router.GET("/signup/github/", Index)
	router.GET("/signup/github/callback/", Index)
	router.GET("/signup/facebook/", Index)
	router.GET("/signup/facebook/callback/", Index)
	router.GET("/signup/google/", Index)
	router.GET("/signup/google/callback/", Index)
	router.GET("/signup/tumblr/", Index)
	router.GET("/signup/tumblr/callback/", Index)

	//login/out
	router.GET("/login/", LoginRender)
	router.POST("/login/", Login)
	router.GET("/login/forgot/", ForgotRender)
	router.POST("/login/forgot/", Index)
	router.GET("/login/reset/", Index)
	router.GET("/login/reset/:email/:token/", Index)
	router.PUT("/login/reset/:email/:token/", Index)
	router.GET("/logout/", Logout)// todo doesn't work

	//social login
	router.GET("/login/twitter/", Index)
	router.GET("/login/twitter/callback/", Index)
	router.GET("/login/github/", Index)
	router.GET("/login/github/callback/", Index)
	router.GET("/login/facebook/", Index)
	router.GET("/login/facebook/callback/", Index)
	router.GET("/login/google/", Index)
	router.GET("/login/google/callback/", Index)
	router.GET("/login/tumblr/", Index)
	router.GET("/login/tumblr/callback/", Index)

	//admin
	//app.all('/admin*', ensureAuthenticated); TODO necessary midleware
	//app.all('/admin*', ensureAdmin);
	router.GET("/admin/", Index)

	//admin > users
	router.GET("/admin/users/", Index)
	router.POST("/admin/users/", Index)
	router.GET("/admin/users/:id/", Index)
	router.PUT("/admin/users/:id/", Index)
	router.PUT("/admin/users/:id/password/", Index)
	router.PUT("/admin/users/:id/role-admin/", Index)
	router.DELETE("/admin/users/:id/role-admin/", Index)
	router.PUT("/admin/users/:id/role-account/", Index)
	router.DELETE("/admin/users/:id/role-account/", Index)
	router.DELETE("/admin/users/:id/", Index)

	//admin > administrators
	router.GET("/admin/administrators/", Index)
	router.POST("/admin/administrators/", Index)
	router.GET("/admin/administrators/:id/", Index)
	router.PUT("/admin/administrators/:id/", Index)
	router.PUT("/admin/administrators/:id/permissions/", Index)
	router.PUT("/admin/administrators/:id/groups/", Index)
	router.PUT("/admin/administrators/:id/user/", Index)
	router.DELETE("/admin/administrators/:id/user/", Index)
	router.DELETE("/admin/administrators/:id/", Index)

	//admin > admin groups
	router.GET("/admin/admin-groups/", Index)
	router.POST("/admin/admin-groups/", Index)
	router.GET("/admin/admin-groups/:id/", Index)
	router.PUT("/admin/admin-groups/:id/", Index)
	router.PUT("/admin/admin-groups/:id/permissions/", Index)
	router.DELETE("/admin/admin-groups/:id/", Index)

	//admin > accounts
	router.GET("/admin/accounts/", Index)
	router.POST("/admin/accounts/", Index)
	router.GET("/admin/accounts/:id/", Index)
	router.PUT("/admin/accounts/:id/", Index)
	router.PUT("/admin/accounts/:id/user/", Index)
	router.DELETE("/admin/accounts/:id/user/", Index)
	router.POST("/admin/accounts/:id/notes/", Index)
	router.POST("/admin/accounts/:id/status/", Index)
	router.DELETE("/admin/accounts/:id/", Index)


	//admin > statuses
	router.GET("/admin/statuses/", Index)
	router.POST("/admin/statuses/", Index)
	router.GET("/admin/statuses/:id/", Index)
	router.PUT("/admin/statuses/:id/", Index)
	router.DELETE("/admin/statuses/:id/", Index)


	//admin > categories
	router.GET("/admin/categories/", Index)
	router.POST("/admin/categories/", Index)
	router.GET("/admin/categories/:id/", Index)
	router.PUT("/admin/categories/:id/", Index)
	router.DELETE("/admin/categories/:id/", Index)

	//admin > search
	router.GET("/admin/search/", Index)

	//account
	//app.all('/account*', ensureAuthenticated); //TODO necessary midleware
	//app.all('/account*', ensureAccount);
	router.GET("/account/", Account)

	//account > verification
	router.GET("/account/verification/", AccountVerification)
	router.POST("/account/verification/", Index)
	router.GET("/account/verification/:token/", Index)

	//account > settings
	router.GET("/account/settings/", AccountSettingsRender)
	router.PUT("/account/settings/", Index)
	router.PUT("/account/settings/identity/", Index)
	router.PUT("/account/settings/password/", Index)

	//account > settings > social
	router.GET("/account/settings/twitter/", Index)
	router.GET("/account/settings/twitter/callback/", Index)
	router.GET("/account/settings/twitter/disconnect/", Index)
	router.GET("/account/settings/github/", Index)
	router.GET("/account/settings/github/callback/", Index)
	router.GET("/account/settings/github/disconnect/", Index)
	router.GET("/account/settings/facebook/", Index)
	router.GET("/account/settings/facebook/callback/", Index)
	router.GET("/account/settings/facebook/disconnect/", Index)
	router.GET("/account/settings/google/", Index)
	router.GET("/account/settings/google/callback/", Index)
	router.GET("/account/settings/google/disconnect/", Index)
	router.GET("/account/settings/tumblr/", Index)
	router.GET("/account/settings/tumblr/callback/", Index)
	router.GET("/account/settings/tumblr/disconnect/", Index)

	//route not found
	router.NoRoute(Status404Render)
}