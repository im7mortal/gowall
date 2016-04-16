package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"strings"
	"regexp"
)

var TemplateStorage map[string]*render.HTML = make(map[string]*render.HTML)

func InitTemplate(base, name string, paths... string) {
	// append base tmpl
	paths = append(paths, "layouts/" + base)
	TemplateStorage[name] = &render.HTML{
		Template: template.Must(template.New(name).ParseFiles(paths...)),
		Name:     base,
	}
}

type Response struct {
	Success bool `json:"success"`
	Errors  []string `json:"errors"`
	ErrFor  map[string]string `json:"errfor"`

	Username    string  `json:"username"`
	Email   string  `json:"email"`
	Password string  `json:"password"`
}

func (r *Response)HasErrors() bool {
	return len(r.ErrFor) != 0 || len(r.Errors) != 0
}

func (r *Response)Fail(c *gin.Context) {
	r.Success = false
	c.JSON(http.StatusOK, r)
}

func (r *Response) Recover(c *gin.Context) {}

func (r *Response) ValidateUsername() {
	r.Username = strings.ToLower(r.Username)
	if len(r.Username) == 0 {
		r.ErrFor["username"] = "required"
	} else {
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\-\_]+$`, r.Username)
		if err != nil {
			println(err.Error())
		}
		if !ok {
			r.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
}

func (r *Response) ValidateEmail() {
	r.Email = strings.ToLower(r.Email)
	if len(r.Email) == 0 {
		r.ErrFor["email"] = "required"
	} else {
		ok, err := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\+]+@[a-zA-Z0-9\-\_\.]+\.[a-zA-Z0-9\-\_]+$`, r.Email)
		if err != nil {
			println(err.Error())
		}
		if !ok {
			r.ErrFor["email"] = `invalid email format`
		}
	}
}

func (r *Response) ValidatePassword() {
	if len(r.Password) == 0 {
		r.ErrFor["password"] = "required"
	} else {
		if len(r.Password) < 8 {
			r.ErrFor["password"] = `too weak password`
		}
	}
}

func (r *Response) CleanErrors() {
	r.Errors = []string{}
	r.ErrFor = make(map[string]string)
}
