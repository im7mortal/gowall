package main

import "os"

const defaultLocalMongoDBUrl  = "mongodb://localhost:27017/gowall"
const defaultPORT  = "3000"
const ROOTGROUP  = "root"


var config struct {
	Port                       string
	CompanyName                string
	ProjectName                string
	SystemEmail                string
	CryptoKey                  string
	RequireAccountVerification bool
	MongoDB                    string
	dbName                     string
	LoginAttempts              LoginAttempts
	SMTP                       SMTP
	Socials                    map[string]OAuth
}

type LoginAttempts struct {
	ForIp         int
	ForIpAndUser  int
	LogExpiration string
}

type SMTP struct {
	From        struct {
				  Name, Address string
			  }
	Credentials struct {
				  User, Password, Host string
				  SSL                  bool
			  }
}

type OAuth struct {
	Key, Secret string
}

func getEnvOrSetDef(envName, defValue string) (val string) {
	val, ok := os.LookupEnv(envName)
	if !ok {
		val = defValue
	}
	return
}

func InitConfig() {

	config.Port = getEnvOrSetDef("PORT", defaultPORT)

	config.CompanyName = "Acme, Inc."
	config.ProjectName = "Gowall"
	config.SystemEmail = "your@email.addy"
	config.CryptoKey = "k3yb0ardc4t"
	config.RequireAccountVerification = true

	config.MongoDB = getEnvOrSetDef(
		"MONGODB_URI",
		getEnvOrSetDef(
			"MONGOLAB_URI",
			getEnvOrSetDef(
				"MONGOHQ_URL",
				defaultLocalMongoDBUrl,
			)))

	if config.dbName == "" {
		config.dbName = getDBName(&config.MongoDB)
	}

	config.LoginAttempts.ForIp = 50
	config.LoginAttempts.ForIpAndUser = 7
	config.LoginAttempts.LogExpiration = "20m"

	config.SMTP.From.Name = getEnvOrSetDef("SMTP_FROM_NAME", config.ProjectName + " Website")
	config.SMTP.From.Address = getEnvOrSetDef("SMTP_FROM_ADDRESS", "your@email.addy")

	//config.SMTP.Credentials.User = getEnvOrSetDef("SMTP_USERNAME", "your@email.addy")
	//config.SMTP.Credentials.Password = getEnvOrSetDef("SMTP_PASSWORD", "bl4rg!")

	config.SMTP.Credentials.User = getEnvOrSetDef("SMTP_USERNAME", "welcome@sturfee.com")
	config.SMTP.Credentials.Password = getEnvOrSetDef("SMTP_PASSWORD", "sturfee_knoxville")
	config.SMTP.Credentials.Host = getEnvOrSetDef("SMTP_HOST", "smtp.gmail.com")
	config.SMTP.Credentials.SSL = true

	// I think it's ok. I use it only for "get". No modifying
	config.Socials = make(map[string]OAuth)

	ins := OAuth{} // todo i hope it's not like JS link

	ins.Key = getEnvOrSetDef("TWITTER_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("TWITTER_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		config.Socials["twitter"] = ins
	}

	ins.Key = getEnvOrSetDef("FACEBOOK_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("FACEBOOK_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		config.Socials["facebook"] = ins
	}

	ins.Key = getEnvOrSetDef("GITHUB_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("GITHUB_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		config.Socials["github"] = ins
	}

	ins.Key = getEnvOrSetDef("GOOGLE_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("GOOGLE_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		config.Socials["google"] = ins
	}

	ins.Key = getEnvOrSetDef("TUMBLR_OAUTH_KEY", "")
	ins.Secret = getEnvOrSetDef("TUMBLR_OAUTH_SECRET", "")
	if len(ins.Key) != 0 {
		config.Socials["tumblr"] = ins
	}

}
