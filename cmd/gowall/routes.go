package main

import (
	"github.com/gin-gonic/gin"
)

func BindRoutes(router *gin.Engine) {

	//front end
	router.GET("/", index)
	router.GET("/about/", about)
	router.GET("/contact/", contactRender)
	router.POST("/contact/", contactSend)

	//sign up
	router.GET("/signup/", signupRender)
	router.POST("/signup/", Signup)

	//social sign up
	router.POST("/signup/social/", SignUpSocial)
	router.GET("/signup/:provider/", SignUpProvider)

	//OAuth callback for /signup/  /login/  /settings/connect/
	router.GET("/provider/:provider/callback/", CompleteOAuth)

	//social login
	router.GET("/provider/:provider/", LoginProvider)

	//login/out
	router.GET("/login/", LoginRender)
	router.POST("/login/", Login)
	router.GET("/login/forgot/", ForgotRender)
	router.POST("/login/forgot/", SendReset)
	router.GET("/login/reset/", ResetRender)
	router.GET("/login/reset/:email/:token/", ResetRender)
	router.PUT("/login/reset/:email/:token/", ResetPassword)
	router.GET("/logout/", Logout)

	//admin
	admin := router.Group("/admin")
	admin.Use(EnsureAuthenticated)
	admin.Use(EnsureAdmin)
	{
		admin.GET("/", renderAdministrator)

		//admin > users
		admin.GET("/users/", renderUsers)
		admin.POST("/users/", createUser)
		admin.GET("/users/:id/", readUser)
		admin.PUT("/users/:id/", changeDataUser)
		admin.PUT("/users/:id/password/", changePasswordUser)
		admin.PUT("/users/:id/role-admin/", linkAdminToUser)
		admin.DELETE("/users/:id/role-admin/", unlinkAdminToUser)
		admin.PUT("/users/:id/role-account/", index)
		admin.DELETE("/users/:id/role-account/", index)
		admin.DELETE("/users/:id/", deleteUser)

		//admin > administrators
		admin.GET("/administrators/", renderAdmin)
		admin.POST("/administrators/", createAdmin)
		admin.GET("/administrators/:id/", readAdmin)
		admin.PUT("/administrators/:id/", updateAdmin)
		admin.PUT("/administrators/:id/permissions/", updatePermissionsAdmin)
		admin.PUT("/administrators/:id/groups/", updateGroupsAdmin) // todo didn't finished
		admin.PUT("/administrators/:id/user/", linkUser)
		admin.DELETE("/administrators/:id/user/", unlinkUser)
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
		admin.PUT("/accounts/:id/", index)
		admin.PUT("/accounts/:id/user/", index)
		admin.DELETE("/accounts/:id/user/", index)
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
		admin.GET("/search/", index) //TODO
	}

	//account
	account := router.Group("/account")
	account.Use(EnsureAuthenticated)
	account.Use(EnsureAccount)
	{
		account.GET("/", AccountRender)

		//account > verification
		account.GET("/verification/", AccountVerificationRender)
		account.POST("/verification/", ResendVerification)
		account.GET("/verification/:token/", Verify)

		//account > settings
		account.GET("/settings/", AccountSettingsRender)
		account.PUT("/settings/", SetSettings)
		account.PUT("/settings/identity/", ChangeIdentity)
		account.PUT("/settings/password/", ChangePassword)

		//account > settings > social
		account.GET("/providerSettings/:provider/", providerSettings)
		account.GET("/providerSettings/:provider/disconnect/", disconnectProvider)
	}

	//route not found
	router.NoRoute(Status404Render)
}
