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
)

const VERSION  = "0.1"

var store sessions.CookieStore

var Router *gin.Engine



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

/*	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          5000 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))*/

	Router.Use(sessions.Sessions("session", store))
	Router.Use(func(c *gin.Context) {

		session := sessions.Default(c)
		oauthMessage, exist := session.Get("oauthMessage").(string)
		session.Delete("oauthMessage")
		session.Save()

		c.Set("oauthMessage", oauthMessage)
		c.Set("oauthMessageExist", exist)
		c.Set("ProjectName", config.ProjectName)
		c.Set("CopyrightYear", "2016") // todo
		c.Set("CopyrightName", config.CompanyName)
		c.Set("CacheBreaker", "br34k-01")
		c.Next()
	})
	Router.Use(IsAuthenticated)
	BindRoutes(Router) // --> cmd/go-getting-started/routers.go

	Router.Run(":" + config.Port)

	/*
	if CONF.DEVELOP {
		router.Run(":8080")
	} else {
		go func() {
			http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
		}()
		errorHTTPS := router.RunTLS(":443", CONF.TLS.CERT, CONF.TLS.KEY)
		if errorHTTPS != nil {
			log.Fatal("HTTPS doesn't work:", errorHTTPS.Error())
		}
	}
	*/
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	//http.Redirect(w, req, "https://" + CONF.ORIGIN + req.RequestURI, http.StatusMovedPermanently)
}

func Logined() bool {
	return true
}