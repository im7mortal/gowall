package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"github.com/im7mortal/gowall/config"
	"regexp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func SignupRender(c *gin.Context) {
/*	if (c.Get("isAuthenticated")) {
		c.Redirect(http.StatusFound, c.Get("defaultReturnUrl"))
	} else {
		_, twitterExist := config.Socials["twitter"]
		_, githubExist := config.Socials["github"]
		_, facebookExist := config.Socials["facebook"]
		_, googleExist := config.Socials["google"]
		_, tumblrExist := config.Socials["tumblr"]
		c.HTML(http.StatusOK, "default.html", gin.H{
			"oauthMessage": "",
			"oauthTwitter": twitterExist,
			"oauthGitHub": githubExist,
			"oauthFacebook": facebookExist,
			"oauthGoogle": googleExist,
			"oauthTumblr": tumblrExist,
		})
	}*/
	//todo
	c.HTML(http.StatusOK, "default.html", gin.H{})
}

func Signup(c *gin.Context) {
	response := Response{} // todo sync.Pool

	defer response.Fail(c)

	username := strings.ToLower(c.Request.FormValue("username"))
	if len(username) == 0 {
		response.ErrFor["username"] = "required"
	} else {
		r, err := regexp.MatchString(`/^[a-zA-Z0-9\-\_]+$/`, username)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
	email := strings.ToLower(c.Request.FormValue("email"))
	if len(email) == 0 {
		response.ErrFor["email"] = "required"
	} else {
		r, err := regexp.MatchString(`/^[a-zA-Z0-9\-\_\.\+]+@[a-zA-Z0-9\-\_\.]+\.[a-zA-Z0-9\-\_]+$/`, email)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["email"] = `invalid email format`
		}
	}
	password := c.Request.FormValue("password")
	if len(password) == 0 {
		response.ErrFor["password"] = "required"
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("test")
	collection := d.C("User")
	//collection.Create(mgo.CollectionInfo{})

	collection.Find(bson.M{"$or": "valera"})


}