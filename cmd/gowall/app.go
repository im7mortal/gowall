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
	"crypto"
	"github.com/xenolf/lego/acme"
	"fmt"
	"log"
	"crypto/rsa"
	"crypto/rand"
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

	initACME()

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



type MyUser struct {
	Email        string
	Registration *acme.RegistrationResource
	key          crypto.PrivateKey
}
func (u MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}
func (u MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}


func initACME()  {

	// Create a user. New accounts need an email and private key to start.
	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		log.Fatal(err)
	}
	myUser := MyUser{
		Email: "im7mortal@gmail.com",
		key: privateKey,
	}

	// A client facilitates communication with the CA server. This CA URL is
	// configured for a local dev instance of Boulder running in Docker in a VM.
	client, err := acme.NewClient("http://192.168.99.100:4000", &myUser, acme.RSA2048)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an http port of 5002 and an tls port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that we still
	// need to proxy challenge traffic to port 5002 and 5001.
	client.SetHTTPAddress(":5002")
	client.SetTLSAddress(":5001")

	// New users will need to register
	reg, err := client.Register()
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	// SAVE THE USER.

	// The client has a URL to the current Let's Encrypt Subscriber
	// Agreement. The user will need to agree to it.
	err = client.AgreeToTOS()
	if err != nil {
		log.Fatal(err)
	}

	// The acme library takes care of completing the challenges to obtain the certificate(s).
	// The domains must resolve to this machine or you have to use the DNS challenge.
	bundle := false
	certificates, failures := client.ObtainCertificate([]string{"go-wall.herokuapp.com"}, bundle, nil)
	if len(failures) > 0 {
		log.Fatal(failures)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)


}
