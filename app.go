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
	"gopkg.in/mgo.v2"
	"github.com/im7mortal/gowall/schemas"
	"gopkg.in/mgo.v2/bson"
)

const VERSION  = "0.1"

var store sessions.CookieStore

var Router *gin.Engine



func init () {
	config.Init()
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
	sess := sessions.Default(c)

	public := sess.Get("public")

	if public != nil && len(public.(string)) > 0 {
		session, err := mgo.Dial("mongodb://localhost:27017")
		defer session.Close()
		if err != nil {
			println(err.Error())
		}
		d := session.DB("test")
		collection := d.C("User")
		us := schemas.User{}
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(public.(string))}).One(&us)
		if err != nil {
			println(err.Error())
		}
		if len(us.Username) > 0 {
			c.Set("Logined", true) // todo what is different between "Logined" and "isAuthenticated"
			c.Set("isAuthenticated", true)
			c.Set("UserName", us.Username)
			c.Set("DefaultReturnUrl", us.DefaultReturnUrl()) // todo
		}
	}
	c.Next()
}

