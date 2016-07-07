package main

import (
	"gopkg.in/gomail.v2"
	"io"
	"html/template"
)

type MailConfig struct {
	From     string
	ReplyTo  string
	To       string
	Subject  string
	TextPath string
	HtmlPath string
	Data     interface{}
}

func (conf *MailConfig)SendMail() (err error) {
	m := gomail.NewMessage()

	m.SetHeader("From", conf.From)
	m.SetHeader("To", "im7mortal@gmail.com") // TODO
	m.SetHeader("Subject", conf.Subject)
	m.SetHeader("ReplyTo", conf.ReplyTo)

	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return template.Must(template.ParseFiles(conf.HtmlPath)).Execute(w, conf.Data)
	})

	d := gomail.NewDialer(config.SMTP.Credentials.Host, 587, config.SMTP.Credentials.User, config.SMTP.Credentials.Password)
	return d.DialAndSend(m)
}

