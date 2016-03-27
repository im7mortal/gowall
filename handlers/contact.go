package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/im7mortal/gowall/config"
	"gopkg.in/gomail.v2"
	"io"
	"html/template"
	"encoding/json"
)

func ContactRender(c *gin.Context) {
	render, _ := TemplateStorage[c.Request.URL.Path]
	render.Data = c.Keys
	c.Render(http.StatusOK, render)
}

func ContactSend(c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)

	//defer response.Recover(c)

	var body struct {
		Name    string  `json:"name"`
		Email   string  `json:"email"`
		Message string  `json:"message"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	if err != nil {
		panic(err.Error())
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
		response.Fail(c)
		return
	}

	m := gomail.NewMessage()

	m.SetHeader("From", config.SMTP.From.Name + " <" + config.SMTP.From.Address + ">")
	m.SetHeader("To", config.SystemEmail)
	m.SetHeader("Subject", config.CompanyName + " contact form")
	m.SetHeader("ReplyTo", body.Email)

	//put in the c.Keys
	c.Set("Name", body.Name)
	c.Set("Email", body.Email)
	c.Set("Message", body.Message)

	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return template.Must(template.ParseFiles("views/contact/email-html.html")).Execute(w, c.Keys)
	})

	d := gomail.NewDialer(config.SMTP.Credentials.Host, 587, config.SMTP.Credentials.User, config.SMTP.Credentials.Password)

	if err := d.DialAndSend(m); err != nil {
		response.Errors = append(response.Errors, "Error Sending: " + err.Error())
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}


