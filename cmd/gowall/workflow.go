package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
)

func handleXHR(c *gin.Context, data []byte) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

type Response struct {
	Success bool `json:"success" bson:"-"`
	Errors  []string `json:"errors" bson:"-"`
	ErrFor  map[string]string `json:"errfor" bson:"-"`
	c  *gin.Context
}

func (r *Response)HasErrors() bool {
	return len(r.ErrFor) != 0 || len(r.Errors) != 0
}

func (r *Response)Fail() {
	r.Success = false
	r.c.JSON(http.StatusOK, r)
}

func (r *Response)BindContext(c *gin.Context) {
	r.c = c
}

func (r *Response) Recover() {}



func (r *Response) CleanErrors() {
	r.Errors = []string{}
	r.ErrFor = make(map[string]string)
}

func (r *Response) DecodeRequest() {
	err := json.NewDecoder(r.c.Request.Body).Decode(r)
	if err != nil {
		panic(err)
	}
	// clean errors from client
	r.CleanErrors()
	return
}

func (r *Response) Finish() {
	r.Success = true
	r.c.JSON(http.StatusOK, r)
}
