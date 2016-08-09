package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func ensureAuthenticated(c *gin.Context) {
	isAuthenticated, _ := c.Get("isAuthenticated")
	if is, ok := isAuthenticated.(bool); ok && is {
		c.Next()
	} else {
		c.Abort()
		session := sessions.Default(c)
		session.Set("returnURL", c.Request.URL.Path)
		session.Save()
		c.Redirect(http.StatusFound, "/login/")
	}
}

func getUser(c *gin.Context) (user *User) {
	if _user, ok := c.Get("User"); ok {
		user, ok = _user.(*User)
		if !ok {
			EXCEPTION("not authorised")
		}
	} else {
		EXCEPTION("not authorised")
	}
	return
}

func getAccount(c *gin.Context) (account *Account) {
	if _account, ok := c.Get("Account"); ok {
		account, ok = _account.(*Account)
		if !ok {
			EXCEPTION("account wasn't found")
		}
	} else {
		EXCEPTION("account wasn't found")
	}
	return
}

func ensureAccount(c *gin.Context) {
	user := getUser(c)
	if ok := user.canPlayRoleOf("account"); ok {
		account := Account{}
		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ACCOUNTS)
		collection.Find(bson.M{"_id": user.Roles.Account}).One(&account)
		c.Set("Account", &account)
		if config.RequireAccountVerification {
			if account.IsVerified != "yes" {
				if yes := rVerificationURL.MatchString(c.Request.URL.Path); !yes {
					c.Redirect(http.StatusFound, "/account/verification/")
					return
				}
			}
		}
		c.Next()
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func getAdmin(c *gin.Context) (admin *Admin) {
	if _admin, ok := c.Get("Admin"); ok {
		admin, ok = _admin.(*Admin)
		if !ok {
			EXCEPTION("user isn't admin")
		}
	} else {
		EXCEPTION("user isn't admin")
	}
	return
}

func ensureAdmin(c *gin.Context) {
	user := getUser(c)
	if ok := user.canPlayRoleOf("admin"); ok {
		admin := Admin{}
		db := getMongoDBInstance()
		defer db.Session.Close()
		collection := db.C(ADMINS)
		collection.Find(bson.M{"_id": user.Roles.Admin}).One(&admin)
		c.Set("Admin", &admin)
		c.Next()
		return
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
			if err != mgo.ErrNotFound {
				EXCEPTION(err)
			}
		}
		if len(us.Username) > 0 {
			c.Set("Logined", true)
			c.Set("isAuthenticated", true)
			c.Set("UserName", us.Username)
			c.Set("User", &us)
			c.Set("DefaultReturnUrl", us.defaultReturnUrl())
		}
	}

	c.Next()
}
