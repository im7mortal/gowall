package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func renderContact(c *gin.Context) {
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func sendContact(c *gin.Context) {
	response := newResponse(c)
	var body struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		EXCEPTION(err.Error())
	}

	if len(body.Name) == 0 {
		response.ErrFor["name"] = "required"
	}
	if len(body.Email) == 0 {
		response.ErrFor["email"] = "required"
	}
	if len(body.Message) == 0 {
		response.ErrFor["message"] = "required"
	}
	if response.HasErrors() {
		response.Fail()
		return
	}

	//put in the c.Keys
	c.Set("Name", body.Name)
	c.Set("Email", body.Email)
	c.Set("Message", body.Message)

	mailConf := MailConfig{}
	mailConf.Data = c.Keys
	mailConf.From = config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">"
	//mailConf.To = config.SystemEmail
	mailConf.To = "im7mortal@gmail.com"
	mailConf.Subject = config.CompanyName + " contact form"
	//mailConf.ReplyTo = body.Email
	mailConf.ReplyTo = "im7mortal@gmail.com"
	mailConf.HtmlPath = "views/contact/email-html.html"

	if err := mailConf.SendMail(); err != nil {
		response.Errors = append(response.Errors, "Email wasn't send. Please try another time or later.")
		response.Fail()
		return
	}

	response.Finish()
}
