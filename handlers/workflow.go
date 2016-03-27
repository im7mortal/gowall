package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/render"
	"html/template"
)

var TemplateStorage map[string]*render.HTML = make(map[string]*render.HTML)

func InitTemplate(base, name string, paths... string) {
	TemplateStorage[name] = &render.HTML{
		Template: template.Must(template.New(name).ParseFiles(paths...)),
		Name:     base,
		Data:     gin.H{},
	}
}

type Response struct {
	Success bool `json:"success"`
	Errors  []string `json:"errors"`
	ErrFor  map[string]string `json:"errfor"`
}

func (r *Response)HasErrors() bool {
	return len(r.ErrFor) != 0 || len(r.Errors) != 0
}

func (r *Response)Fail(c *gin.Context) {
	r.Success = false
	c.JSON(http.StatusOK, r)
}