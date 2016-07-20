/** Response object was a struct but i change it to map because of in drywall often some
	data was in response. I mean every response has mandatory fields but sometimes also
	custom data like "newUsername". I don't like it. I think in response type has to be
	field data map[string]interface{}. But my principe is "don't touch client code". And I
	decided do response object as map. It will be flexible for users from node.js.

	type Response struct {
	Success bool `json:"success" bson:"-"`
	Errors  []string `json:"errors" bson:"-"`
	ErrFor  map[string]string `json:"errfor" bson:"-"`
	c  *gin.Context
	Data map[string]interface{} `json:"data" bson:"-"`
}

	I thought leave struct realization.
	I thought do setters and getters. But it will complicate code. Go is not about
	complexity.
*/

package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleXHR(c *gin.Context, data []byte) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

type Response struct {
	Success bool              `json:"success" bson:"-"`
	Errors  []string          `json:"errors" bson:"-"`
	ErrFor  map[string]string `json:"errfor" bson:"-"`
	c       *gin.Context
	Data    map[string]interface{} `json:"data" bson:"-"`
}

func (r *Response) HasErrors() bool {
	return len(r.ErrFor) != 0 || len(r.Errors) != 0
}

func (r *Response) Fail() {
	r.Success = false
	r.Response()
}

func (r *Response) Init(c *gin.Context) {
	r.c = c
	r.Data = map[string]interface{}{}
	r.ErrFor = map[string]string{}
}

func getResponseObj(c *gin.Context) (r *Response) {
	r = &Response{}
	r.c = c
	r.Data = map[string]interface{}{}
	r.ErrFor = map[string]string{}
	return
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
	r.Errors = []string{}
	r.ErrFor = make(map[string]string)
	return
}

func (r *Response) Finish() {
	r.Success = true
	r.Response()
}

func (r *Response) Response() {
	r.Data["success"] = r.Success
	r.Data["errfor"] = r.ErrFor
	r.Data["errors"] = r.Errors
	r.c.JSON(http.StatusOK, r.Data)
}
