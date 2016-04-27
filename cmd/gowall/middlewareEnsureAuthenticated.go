package main

import (
"github.com/gin-gonic/contrib/sessions"
"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
"gopkg.in/mgo.v2/bson"
)

const USERS  = "users"
const LOGINATTEMPTS  = "loginattempts"
const ACCOUNTS  = "accounts"
const ADMINGROUPS  = "admingroups"
const CATEGORIES  = "categories"
const STATUSES = "status"
const ADMINS  = "admins"

func EnsureAuthenticated(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		c.Next()
	} else {
		session := sessions.Default(c)
		session.Set("returnURL", c.Request.URL.Path)
		session.Save()
		c.Redirect(http.StatusFound, "/login/")
	}
}

func getUser(c *gin.Context) (user *User, ok bool) {
	if _user, _ok := c.Get("User"); _ok {
		user, ok = _user.(*User)
	}
	return
}

func getAccount(c *gin.Context) (account *Account, ok bool) {
	if _account, _ok := c.Get("Account"); _ok {
		account, ok = _account.(*Account)
	}
	return
}

func EnsureAccount(c *gin.Context) {
	if user, ok := getUser(c); ok {
		if ok = user.CanPlayRoleOf("account"); ok {
			account := Account{}
			db := getMongoDBInstance()
			defer db.Session.Close()
			collection := db.C(ACCOUNTS)
			collection.Find(bson.M{"_id": user.Roles.Account}).One(&account)
			c.Set("Account", &account)
			if config.RequireAccountVerification {
				if account.IsVerified != "yes" {
					r, _ := regexp.MatchString(`^\/account\/verification\/`, c.Request.URL.Path)
					if !r {
						c.Redirect(http.StatusFound, "/account/verification/")
					}
				}
			}
			c.Next()
			return
		}
	}
	c.Redirect(http.StatusFound, "/")
}

func getAdmin(c *gin.Context) (account *Admin, ok bool) {
	if _account, _ok := c.Get("Admin"); _ok {
		account, ok = _account.(*Admin)
	}
	return
}

func EnsureAdmin(c *gin.Context) {
	if user, ok := getUser(c); ok {
		if ok = user.CanPlayRoleOf("admin"); ok {
			admin := Admin{}
			db := getMongoDBInstance()
			defer db.Session.Close()
			collection := db.C(ADMINS)
			collection.Find(bson.M{"_id": user.Roles.Admin}).One(&admin)
			c.Set("Admin", &admin)
			c.Next()
			return
		}
	}
	c.Redirect(http.StatusFound, "/")
}

func IsAuthenticated(c *gin.Context) {
	c.Set("isAuthenticated", false)

	session := sessions.Default(c)

	public := session.Get("public")
	public_, ok := public.(string)
	if ok && len(public_) > 0 {
		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(USERS)
		us := User{}
		err := collection.Find(bson.M{"_id": bson.ObjectIdHex(public_)}).One(&us)
		if err != nil {
			println(err.Error())
		}
		if len(us.Username) > 0 {
			c.Set("Logined", true) // todo what is different between "Logined" and "isAuthenticated"
			c.Set("isAuthenticated", true)
			c.Set("UserName", us.Username)
			c.Set("User", &us)
			c.Set("DefaultReturnUrl", us.DefaultReturnUrl())
		}
	}

	c.Next()
}
