package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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