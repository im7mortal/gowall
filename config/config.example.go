package config

import (
	"os"
)

var Port = getEnvOrSetDef("PORT", "3000")

const (
	CompanyName string = "Acme, Inc."
	ProjectName string = "Gowall"
	SystemEmail string = "your@email.addy"
	CryptoKey string = "k3yb0ardc4t"
	RequireAccountVerification bool = true
)

var MongoDB = getEnvOrSetDef(
	"process.env.MONGOLAB_URI",
	getEnvOrSetDef("process.env.MONGOHQ_URL", "mongodb://localhost:27017/drywall"))

var LoginAttempts struct {
	ForIp         int
	ForIpAndUser  int
	LogExpiration string
}

var SMTP struct {
	From        struct {
				  Name, Address string
			  }
	Credentials struct {
				  User, Password, Host string
				  SSL                  bool
			  }
}

type OAuth struct{
	Key, Secret string
}

// I think it's ok. I use it only for "get". No modifying
var Socials = make(map[string]OAuth)

func getEnvOrSetDef(envName, defValue string) (val string) {
	val, ok := os.LookupEnv(envName)
	if !ok {
		val = defValue
	}
	return
}

func Init() {

	LoginAttempts.ForIp = 50
	LoginAttempts.ForIpAndUser = 7
	LoginAttempts.LogExpiration = "20m"

	SMTP.From.Name = getEnvOrSetDef("SMTP_FROM_NAME", ProjectName + " Website")
	SMTP.From.Address = getEnvOrSetDef("SMTP_FROM_ADDRESS", "your@email.addy")

	SMTP.Credentials.User = getEnvOrSetDef("SMTP_USERNAME", "your@email.addy")
	SMTP.Credentials.Password = getEnvOrSetDef("SMTP_PASSWORD", "bl4rg!")
	SMTP.Credentials.Host = getEnvOrSetDef("SMTP_HOST", "smtp.gmail.com")
	SMTP.Credentials.SSL = true

	ins := OAuth{} // todo i hope it's not like JS link

	ins.Key = getEnvOrSetDef("TWITTER_OAUTH_KEY", "11")
	ins.Secret = getEnvOrSetDef("TWITTER_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		Socials["twitter"] = ins
	}

	ins.Key = getEnvOrSetDef("FACEBOOK_OAUTH_KEY", "11")
	ins.Secret = getEnvOrSetDef("FACEBOOK_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		Socials["facebook"] = ins
	}

	ins.Key = getEnvOrSetDef("GITHUB_OAUTH_KEY", "11")
	ins.Secret = getEnvOrSetDef("GITHUB_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		Socials["github"] = ins
	}

	ins.Key = getEnvOrSetDef("GOOGLE_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("GOOGLE_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		Socials["google"] = ins
	}

	ins.Key = getEnvOrSetDef("TUMBLR_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("TUMBLR_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		Socials["tumblr"] = ins
	}

}