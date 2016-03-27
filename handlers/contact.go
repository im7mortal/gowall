package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/im7mortal/gowall/config"
	"gopkg.in/gomail.v2"
)

func ContactRender(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func ContactSend(c *gin.Context) {
	response := Response{} // todo sync.Pool

	defer response.Fail(c)

	name := c.Request.FormValue("name")
	if len(name) == 0 {
		response.ErrFor["name"] = "required"
	}
	email := c.Request.FormValue("email")
	if len(email) == 0 {
		response.ErrFor["email"] = "required"
	}
	message := c.Request.FormValue("message")
	if len(message) == 0 {
		response.ErrFor["message"] = "required"
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}



	m := gomail.NewMessage()

	m.SetHeader("From", config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">")
	m.SetHeader("To", config.SystemEmail)
	m.SetHeader("Subject", config.CompanyName + " contact form")
	m.SetHeader("ReplyTo", email)

	/*
		m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return emailTmpl.Execute(w, data)
	})

	*/


	/*d := gomail.NewDialer(config.SMTP.From.Name, 587, CONF.MAIL.USERNAME, CONF.MAIL.PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		response.Errors = append(response.Errors, "Error Sending: " + err.Error())
		response.Fail(c)
		return
	}*/


	response.Success = true
	c.JSON(http.StatusOK, response)
}


