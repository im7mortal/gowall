package main

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/sessions"
	//"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	//"github.com/itsjamie/gin-cors"
	//"time"
	//"log"
	"net/http"
	"time"
)

const VERSION  = "0.1"

var store sessions.CookieStore

var Router *gin.Engine

var year int

func init () {
	InitConfig()
	store = sessions.NewCookieStore([]byte("MFDQmJQ4TF"))
	store.Options(sessions.Options{
		Path: "/",
		MaxAge: 60 * 60 * 6,
		//Secure: !CONF.DEVELOP,  // TODO https sudo run hack for test without https on local machine
		HttpOnly: true,
	})

}

func main() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.New()
	Router.Use(gzip.Gzip(gzip.DefaultCompression))

	Router.StaticFile("/favicon.ico", "public/favicon.ico")
	Router.Static("/public", "public")
	Router.Static("/vendor", "vendor") // todo not good. conflict with go project structure

	// templates
	Router.HTMLRender = initTemplates(Router)

	Router.Use(gin.Logger())
	Router.Use(checkRecover)

	Router.Use(sessions.Sessions("session", store))

	//refresh year every minute
	go func() {
		for ;; {
			year, _, _ = time.Now().Date()
			time.Sleep(time.Minute)
		}
	} ()

	Router.Use(func(c *gin.Context) {

		session := sessions.Default(c)
		oauthMessage, exist := session.Get("oauthMessage").(string)
		session.Delete("oauthMessage")
		session.Save()

		c.Set("oauthMessage", oauthMessage)
		c.Set("oauthMessageExist", exist)
		c.Set("ProjectName", config.ProjectName)
		c.Set("CopyrightYear", year)
		c.Set("CopyrightName", config.CompanyName)
		c.Set("CacheBreaker", "br34k-01")
		c.Next()
	})
	Router.Use(IsAuthenticated)
	BindRoutes(Router) // --> cmd/go-getting-started/routers.go

	Router.Run(":" + config.Port)

	// https
	// put path to cert instead of CONF.TLS.CERT
	// put path to key instead of CONF.TLS.KEY
	/*
	go func() {
			http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
		}()
		errorHTTPS := router.RunTLS(":443", CONF.TLS.CERT, CONF.TLS.KEY)
		if errorHTTPS != nil {
			log.Fatal("HTTPS doesn't work:", errorHTTPS.Error())
		}
	*/
}

// force redirect to https from http
// necessary only if you use https directly
// put your domain name instead of CONF.ORIGIN
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	//http.Redirect(w, req, "https://" + CONF.ORIGIN + req.RequestURI, http.StatusMovedPermanently)
}
