package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"net/http"
)

func init() {
	gothic.Store = store
}

func startOAuth(c *gin.Context) {
	// don't like that hack
	// gothic was written for another path
	// I just put provider query
	provider := c.Param("provider")
	c.Request.URL.RawQuery += "provider=" + provider

	checkProvider(provider, c.Request.Host)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func completeOAuth(c *gin.Context) {
	// gothic was written for another path
	// i just put provider query
	provider := c.Param("provider")
	c.Request.URL.RawQuery += "&provider=" + provider

	userGoth, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	session := sessions.Default(c)
	action, ok := session.Get("action").(string)
	if !ok {
		EXCEPTION("OAuth action isn't defined")
	}
	session.Delete("action")
	session.Save()

	if err != nil {
		c.HTML(http.StatusOK, action, c.Keys)
		return
	}

	switch action {
	case "/signup/":
		signupProvider(c, userGoth)
		return
	case "/login/":
		loginProvider(c, userGoth)
		return
	case "/account/settings/":
		settingsProvider(c, userGoth)
		return
	default:
		EXCEPTION("OAuth action isn't defined")
	}
}

func checkProvider(provider, hostname string) {
	_, err := goth.GetProvider(provider)
	if err != nil {
		callbackURL := "http://" + hostname + "/provider/" + provider + "/callback"
		switch provider {
		case "facebook":
			goth.UseProviders(
				facebook.New(config.Socials[provider].Key, config.Socials[provider].Secret, callbackURL),
			)
			return
		case "github":
			goth.UseProviders(
				github.New(config.Socials[provider].Key, config.Socials[provider].Secret, callbackURL),
			)
			return
		default:
			EXCEPTION("provider doesn't exist")
		}
	}
}

func (user *User) updateProvider(socialProfile goth.User) {
	switch socialProfile.Provider {
	case "facebook":
		user.Facebook = vendorOauth{}
		user.Facebook.ID = socialProfile.UserID
		return
	case "github":
		user.Github = vendorOauth{}
		user.Github.ID = socialProfile.UserID
		return
	default:
		EXCEPTION("provider doesn't exist")
	}
}

func (user *User) disconnectProviderDB(provider string) {
	switch provider {
	case "facebook":
		user.Facebook.ID = ""
		return
	case "github":
		user.Github.ID = ""
		return
	default:
		EXCEPTION("provider doesn't exist")
	}
}

func injectSocials(c *gin.Context) {

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
}

func doUserHasSocials(c *gin.Context, user *User) {
	if len(user.Facebook.ID) != 0 {
		c.Set("oauthFacebookActive", true)
	}
	if len(user.Twitter.ID) != 0 {
		c.Set("oauthTwitterActive", true)
	}
	if len(user.Github.ID) != 0 {
		c.Set("oauthGitHubActive", true)
	}
	if len(user.Google.ID) != 0 {
		c.Set("oauthGoogleActive", true)
	}
	if len(user.Tumblr.ID) != 0 {
		c.Set("oauthTumblrActive", true)
	}
}
