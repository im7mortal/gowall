package main

import (
	"github.com/gin-gonic/gin"
)

func bindRoutes(router *gin.Engine) {

	//front end
	router.GET("/", index)
	router.GET("/about/", about)
	router.GET("/contact/", renderContact)
	router.POST("/contact/", sendContact)

	//sign up
	router.GET("/signup/", renderSignup)
	router.POST("/signup/", signup)

	//social sign up
	router.POST("/signup/social/", socialSignup)
	router.GET("/signup/:provider/", providerSignup)

	//OAuth callback for /signup/  /login/  /settings/connect/
	router.GET("/provider/:provider/callback/", completeOAuth)

	//social login
	router.GET("/provider/:provider/", providerLogin)

	//login/out
	router.GET("/login/", renderLogin)
	router.POST("/login/", login)
	router.GET("/login/forgot/", renderForgot)
	router.POST("/login/forgot/", sendReset)
	router.GET("/login/reset/", renderReset)
	router.GET("/login/reset/:email/:token/", renderReset)
	router.PUT("/login/reset/:email/:token/", resetPassword)
	router.GET("/logout/", logout)

	//admin
	admin := router.Group("/admin")
	admin.Use(ensureAuthenticated)
	admin.Use(ensureAdmin)
	{
		admin.GET("/", renderAdministrator)

		//admin > users
		admin.GET("/users/", renderUsers)
		admin.POST("/users/", createUser)
		admin.GET("/users/:id/", readUser)
		admin.PUT("/users/:id/", changeDataUser)
		admin.PUT("/users/:id/password/", changePasswordUser)
		admin.PUT("/users/:id/role-admin/", linkAdminToUser)
		admin.DELETE("/users/:id/role-admin/", unlinkAdminFromUser)
		admin.PUT("/users/:id/role-account/", linkAccountToUser)
		admin.DELETE("/users/:id/role-account/", unlinkAccountFromUser)
		admin.DELETE("/users/:id/", deleteUser)

		//admin > administrators
		admin.GET("/administrators/", renderAdmins)
		admin.POST("/administrators/", createAdmin)
		admin.GET("/administrators/:id/", readAdmin)
		admin.PUT("/administrators/:id/", updateAdmin)
		admin.PUT("/administrators/:id/permissions/", updatePermissionsAdmin)
		admin.PUT("/administrators/:id/groups/", updateGroupsAdmin)
		admin.PUT("/administrators/:id/user/", linkUserToAdmin)
		admin.DELETE("/administrators/:id/user/", unlinkUserFromAdmin)
		admin.DELETE("/administrators/:id/", deleteAdmin)

		//admin > admin groups
		admin.GET("/admin-groups/", renderAdminGroups)
		admin.POST("/admin-groups/", createAdminGroup)
		admin.GET("/admin-groups/:id/", readAdminGroup)
		admin.PUT("/admin-groups/:id/", updateAdminGroup)
		admin.PUT("/admin-groups/:id/permissions/", updatePermissionsAdminGroup)
		admin.DELETE("/admin-groups/:id/", deleteAdminGroup)

		//admin > accounts
		admin.GET("/accounts/", renderAccounts)
		admin.POST("/accounts/", createAccount)
		admin.GET("/accounts/:id/", readAccount)
		admin.PUT("/accounts/:id/", updateAccount)
		admin.PUT("/accounts/:id/user/", linkUserToAccount)
		admin.DELETE("/accounts/:id/user/", unlinkUserFromAccount)
		admin.POST("/accounts/:id/notes/", newNote)
		admin.POST("/accounts/:id/status/", newStatus)
		admin.DELETE("/accounts/:id/", deleteAccount)

		//admin > statuses
		admin.GET("/statuses/", renderStatuses)
		admin.POST("/statuses/", createStatus)
		admin.GET("/statuses/:id/", readStatus)
		admin.PUT("/statuses/:id/", updateStatus)
		admin.DELETE("/statuses/:id/", deleteStatus)

		//admin > categories
		admin.GET("/categories/", renderCategories)
		admin.POST("/categories/", createCategory)
		admin.GET("/categories/:id/", renderCategory)
		admin.PUT("/categories/:id/", updateCategory)
		admin.DELETE("/categories/:id/", deleteCategory)

		//admin > search
		admin.GET("/search/", searchResult)
	}

	//account
	account := router.Group("/account")
	account.Use(ensureAuthenticated)
	account.Use(ensureAccount)
	{
		account.GET("/", renderAccount)

		//account > verification
		account.GET("/verification/", renderAccountVerification)
		account.POST("/verification/", resendVerification)
		account.GET("/verification/:token/", verify)

		//account > settings
		account.GET("/settings/", renderAccountSettings)
		account.PUT("/settings/", setSettings)
		account.PUT("/settings/identity/", changeIdentity)
		account.PUT("/settings/password/", changePassword)

		//account > settings > social
		account.GET("/providerSettings/:provider/", settingsProvider_)
		account.GET("/providerSettings/:provider/disconnect/", disconnectProvider)
	}

	//route not found
	router.NoRoute(renderStatus404)
}
