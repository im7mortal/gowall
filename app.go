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



	LoadTemplates()


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
	Router.Use(IsAuthenticated)
	Router.Use(func(c *gin.Context) {
		c.Set("ProjectName", config.ProjectName)
		c.Set("CopyrightYear", "2016") // todo
		c.Set("CopyrightName", config.CompanyName)
		c.Set("CacheBreaker", "br34k-01")
		c.Next()
	})

	Router.Static("/public", "public")
	Router.Static("/vendor", "vendor")
	Router.Static("/views", "views") // todo not good

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


func IsAuthenticated(c *gin.Context) {
	/*path := c.Request.URL.Path
	if (len(path) > 5 && path[:6] == "/auth/") {
		c.Next()
		return
	}
	session := getSession(c)
	// reset cookies for mowing expiration time
	// TODO  don't sure it's good
	datr, err := c.Request.Cookie("datr")
	if err == nil {
		setCookieForJS("datr", datr.Value, c)
		session.Set("session", session.Get("session"))
		session.Save()
	}*/
	c.Next()
}

