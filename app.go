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
	"html/template"
)

const VERSION  = "0.1"

var store sessions.CookieStore

var Router *gin.Engine



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
	Router = gin.Default()
	Router.Use(gzip.Gzip(gzip.DefaultCompression))





	indexTemplate := template.Must(template.ParseFiles("layouts/default.html"))
	tm1 = template.Must(template.Must(indexTemplate.Clone()).ParseFiles("body.html"))
	tm2 = template.Must(template.Must(indexTemplate.Clone()).ParseFiles("body2.html"))

	ty, _ := template.ParseFiles("layouts/default.html", "body.html", "body2.html")


	Router.SetHTMLTemplate(ty)




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

	Router.LoadHTMLFiles("layouts/default.html", "body.html", "body2.html")

	//router.LoadHTMLFiles(CONF.PATH + "/t_web/layouts/module1/module1.html")
	//router.Use(static.Serve("/", static.LocalFile(CONF.PATH+"/t_knoxville", true)))



	Router.Static("/public", "public")

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

var trigger bool
var tm1 *template.Template
var tm2 *template.Template


