package main

// create collections if not exist and ensure indexes

const USERS = "users"
const LOGINATTEMPTS = "loginattempts"
const ACCOUNTS = "accounts"
const ADMINGROUPS = "admingroups"
const CATEGORIES = "categories"
const STATUSES = "status"
const ADMINS = "admins"

func init() {

	db := getMongoDBInstance()
	defer db.Session.Close()

	c := db.C(USERS)
	c.EnsureIndex(UserUniqueIndex)
	c.EnsureIndex(UserIndex)

	c = db.C(LOGINATTEMPTS)
	c.EnsureIndex(LoginAttemptsIndex)

	c = db.C(ACCOUNTS)
	c.EnsureIndex(AccountIndex)

	c = db.C(ADMINGROUPS)
	c.EnsureIndex(AdminGroupIndex)

	c = db.C(CATEGORIES)
	c.EnsureIndex(CategoryIndex)

	c = db.C(STATUSES)
	c.EnsureIndex(StatusesIndex)

	c = db.C(ADMINS)
	c.EnsureIndex(AdminsIndex)

}
