package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
"github.com/im7mortal/gowall/config"
)

func LoginRender(c *gin.Context) {
	_, isAuthenticated := c.Get("isAuthenticated") // non standard way. If exist it isAuthenticated
	if isAuthenticated {
		defaultReturnUrl, _ := c.Get("defaultReturnUrl")
		c.Redirect(http.StatusFound, defaultReturnUrl.(string))
	} else {
		render, _ := TemplateStorage[c.Request.URL.Path]

		_, oauthTwitter := config.Socials["twitter"]
		_, oauthGitHub := config.Socials["github"]
		_, oauthFacebook := config.Socials["facebook"]
		_, oauthGoogle := config.Socials["google"]
		_, oauthTumblr := config.Socials["tumblr"]

		c.Set("oauth", oauthTwitter || oauthGitHub || oauthFacebook || oauthGoogle || oauthTumblr)
		c.Set("oauthTwitter", oauthTwitter)
		c.Set("oauthGitHub", oauthGitHub)
		c.Set("oauthFacebook", oauthFacebook)
		c.Set("oauthGoogle", oauthGoogle)
		c.Set("oauthTumblr", oauthTumblr)
		c.Set("oauthMessage", "")

		render.Data = c.Keys
		c.Render(http.StatusOK, render)
	}
}
