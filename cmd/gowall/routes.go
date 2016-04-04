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
	admin := router.Group("/admin")
	admin.Use(EnsureAuthenticated)
	//admin.Use(ensureAdmin)
	{
		admin.GET("/admin/", Index)

		//admin > users
		admin.GET("/users/", Index)
		admin.POST("/users/", Index)
		admin.GET("/users/:id/", Index)
		admin.PUT("/users/:id/", Index)
		admin.PUT("/users/:id/password/", Index)
		admin.PUT("/users/:id/role-admin/", Index)
		admin.DELETE("/users/:id/role-admin/", Index)
		admin.PUT("/users/:id/role-account/", Index)
		admin.DELETE("/users/:id/role-account/", Index)
		admin.DELETE("/users/:id/", Index)

		//admin > administrators
		admin.GET("/administrators/", Index)
		admin.POST("/administrators/", Index)
		admin.GET("/administrators/:id/", Index)
		admin.PUT("/administrators/:id/", Index)
		admin.PUT("/administrators/:id/permissions/", Index)
		admin.PUT("/administrators/:id/groups/", Index)
		admin.PUT("/administrators/:id/user/", Index)
		admin.DELETE("/administrators/:id/user/", Index)
		admin.DELETE("/administrators/:id/", Index)

		//admin > admin groups
		admin.GET("/admin-groups/", Index)
		admin.POST("/admin-groups/", Index)
		admin.GET("/admin-groups/:id/", Index)
		admin.PUT("/admin-groups/:id/", Index)
		admin.PUT("/admin-groups/:id/permissions/", Index)
		admin.DELETE("/admin-groups/:id/", Index)

		//admin > accounts
		admin.GET("/accounts/", Index)
		admin.POST("/accounts/", Index)
		admin.GET("/accounts/:id/", Index)
		admin.PUT("/accounts/:id/", Index)
		admin.PUT("/accounts/:id/user/", Index)
		admin.DELETE("/accounts/:id/user/", Index)
		admin.POST("/accounts/:id/notes/", Index)
		admin.POST("/accounts/:id/status/", Index)
		admin.DELETE("/accounts/:id/", Index)

		//admin > statuses
		admin.GET("/statuses/", Index)
		admin.POST("/statuses/", Index)
		admin.GET("/statuses/:id/", Index)
		admin.PUT("/statuses/:id/", Index)
		admin.DELETE("/statuses/:id/", Index)

		//admin > categories
		admin.GET("/categories/", Index)
		admin.POST("/categories/", Index)
		admin.GET("/categories/:id/", Index)
		admin.PUT("/categories/:id/", Index)
		admin.DELETE("/categories/:id/", Index)

		//admin > search
		admin.GET("/search/", Index)
	}

	//account
	account := router.Group("/account")
	account.Use(EnsureAuthenticated)
	//account.Use(ensureAccount)
	{
		account.GET("/", AccountRender)

		//account > verification
		account.GET("/verification/", AccountVerification)
		account.POST("/verification/", Index)
		account.GET("/verification/:token/", Index)

		//account > settings
		account.GET("/settings/", AccountSettingsRender)
		account.PUT("/settings/", Index)
		account.PUT("/settings/identity/", Index)
		account.PUT("/settings/password/", Index)

		//account > settings > social
		account.GET("/settings/twitter/", Index)
		account.GET("/settings/twitter/callback/", Index)
		account.GET("/settings/twitter/disconnect/", Index)
		account.GET("/settings/github/", Index)
		account.GET("/settings/github/callback/", Index)
		account.GET("/settings/github/disconnect/", Index)
		account.GET("/settings/facebook/", Index)
		account.GET("/settings/facebook/callback/", Index)
		account.GET("/settings/facebook/disconnect/", Index)
		account.GET("/settings/google/", Index)
		account.GET("/settings/google/callback/", Index)
		account.GET("/settings/google/disconnect/", Index)
		account.GET("/settings/tumblr/", Index)
		account.GET("/settings/tumblr/callback/", Index)
		account.GET("/settings/tumblr/disconnect/", Index)
	}

	//route not found
	router.NoRoute(Status404Render)
}