package main

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/im7mortal/gowall/config"
	"github.com/gin-gonic/contrib/sessions"
	//"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	//"github.com/itsjamie/gin-cors"
	//"time"
	//"log"
	"net/http"
)

var store sessions.CookieStore

func init () {
	config.Init()
/*	store = sessions.NewCookieStore([]byte("MFDQmJQ4TFHVad3dEddQV8QkSUFBUkFi1CQkCRHVad3dEdUFB03HVad3dEddEdi1C"))
	store.Options(sessions.Options{
		Path: "/",
		MaxAge: 60 * 60 * 6,
		Secure: !CONF.DEVELOP,  // TODO https sudo run hack for test without https on local machine
		HttpOnly: true,
	})*/
}

func main() {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

/*	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          5000 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))*/

	router.Use(sessions.Sessions("session", store))

	router.LoadHTMLGlob("layouts/default.html")

	//router.LoadHTMLFiles(CONF.PATH + "/t_web/layouts/module1/module1.html")
	//router.Use(static.Serve("/", static.LocalFile(CONF.PATH+"/t_knoxville", true)))

	BindRoutes(router) // --> cmd/go-getting-started/routers.go

	router.Static("/public", "public")


	router.Run(":" + config.Port)

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

func Smock(c *gin.Context) {

	c.HTML(http.StatusOK, "default.html", gin.H{
		"title": "Main website",
		"pipeline": true,
	})
}