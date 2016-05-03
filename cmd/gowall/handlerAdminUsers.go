package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"html/template"
	"strings"

	"gopkg.in/mgo.v2"
	"net/url"
	"golang.org/x/crypto/bcrypt"
)

func AdminUsersRender(c *gin.Context) {
	query := bson.M{}

	username, ok := c.GetQuery("username")
	if ok && len(username) != 0 {
		query["username"] = bson.RegEx{
			Pattern: `^.*?` + username + `.*$`,
			Options: "i",
		}
	}

	isActive, ok := c.GetQuery("isActive")
	if ok && len(isActive) != 0 {
		query["isActive"] = isActive
	}

	roles, ok := c.GetQuery("roles")
	if ok && len(roles) != 0 {
		// roles.admin or roles.account
		query["roles." + roles] = bson.M{
			"$exists": true,
		}
	}

	type _user struct {
		ID bson.ObjectId `bson:"_id" json:"_id"`
		Username string `bson:"username" json:"username"`
		IsActive string `bson:"isActive" json:"isActive"`
		Email string `bson:"email" json:"email"`
	}
	var results []_user

	db := getMongoDBInstance()
	defer db.Session.Close()
	// TODO keys
	collection := db.C(USERS)

	Result := getData(c, collection.Find(query), &results)



	Results, _ := json.Marshal(Result)

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", Results)
		return
	}

	c.Set("Results", template.JS(string(Results)))
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

func getData (c *gin.Context, query *mgo.Query, results interface{}) (data gin.H) {
	limitS := c.DefaultQuery("limit", "20")
	limit_, _ := strconv.ParseInt(limitS, 0, 0)
	limit := int(limit_)
	if limit > 100 {
		limit = 100
	}

	pageS := c.DefaultQuery("page", "0")
	page_, _ := strconv.ParseInt(pageS, 0, 0)
	page := int(page_)
	sort := c.DefaultQuery("sort", "_id")

	count, _ := query.Count()
	query.Skip(page * limit).Sort(sort).Limit(limit).All(results)

	page += 1
	count_ := page * limit
	pages := gin.H{
		"current": page,
		"prev": page - 1,
		"hasPrev": page - 1 != 0,
		"next": page + 1,
		"hasNext": float64(count) / float64(count_) > 1,
		"total": count,
	}

	end := count_
	if count_ > count {
		end = count
	}

	items := gin.H{
		"begin": (page - 1) * limit,
		"end": end,
		"total": count,
	}

	return gin.H{
		"data": results,
		"pages": pages,
		"items": items,
	}
}


func CreateUser(c *gin.Context) {
	response := Response{} // todo sync.Pool
	defer response.Recover(c)

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&response)
	if err != nil {
		panic(err)
		return
	}
	// clean errors from client
	response.CleanErrors()

	// validate
	response.ValidateUsername()

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.Find(bson.M{"username": response.Username}).One(&user)
	if err != nil {
		println(err.Error())
	}

	// duplicateUsernameCheck
	if len(user.Username) != 0 {
		if user.Username == response.Username {
			response.Errors = append(response.Errors, "That username is already taken.")
		}
	}
	if response.HasErrors() {
		response.Fail(c)
		return
	}

	// createUser

	user.ID = bson.NewObjectId()
	user.Username = response.Username
	user.Search = []string{response.Username}

	err = collection.Insert(user)
	if err != nil {
		panic(err)
		return
	}
	response.Success = true
	c.JSON(http.StatusOK, response)
}

func UsersRender(c *gin.Context) {

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	json, _ := json.Marshal(gin.H{
		"_id": user.ID.Hex(),
		"username": user.Username,
		"email": user.Email,
		"isActive": user.IsActive,
	})

	if XHR(c) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Data(http.StatusOK, "application/json; charset=utf-8", json)
		return
	}

	c.Set("Record", template.JS(url.QueryEscape(string(json))))
	c.HTML(http.StatusOK, "/admin/users/details/", c.Keys)
}

func DeleteUser(c *gin.Context) {
	admin := getAdmin(c)
	user := getUser(c)

	response := Response{} // todo sync.Pool
	defer response.Recover(c)

	// validate
	ok := admin.IsMemberOf("root")
	if !ok {
		response.Errors = append(response.Errors, "You may not delete users.")
		response.Fail(c)
		return
	}

	deleteID := c.Param("id")

	if deleteID == user.ID.Hex() {
		response.Errors = append(response.Errors, "You may not delete yourself from user.")
		response.Fail(c)
		return
	}

	// deleteUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	err := collection.RemoveId(bson.ObjectIdHex(deleteID))
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

func ChangeUserPassword (c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)
	var body struct {
		Confirm   string  `json:"confirm"`
		Password string  `json:"newPassword"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	// validate
	if len(body.Password) == 0 {
		response.ErrFor["newPassword"] = "required"
	}
	if len(body.Confirm) == 0 {
		response.ErrFor["confirm"] = "required"
	} else if body.Password != body.Confirm {
		response.Errors = append(response.Errors, "Passwords do not match.")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	// patchUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail(c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Errors = append(response.Errors, err.Error()) // TODO don't like that this error goes to client
		response.Fail(c)
		return
	}

	user.Password = string(hashedPassword)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}

/*
func ChangeUserData (c *gin.Context) {
	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)
	var body struct {
		Username    string  `json:"username"`
		Email   string  `json:"email"`
		IsActive   string  `json:"isActive"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)


	username := strings.ToLower(body.Username)
	if len(username) == 0 {
		response.ErrFor["username"] = "required"
	} else {
		r, err := regexp.MatchString(`^[a-zA-Z0-9\-\_]+$`, username)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["username"] = `only use letters, numbers, \'-\', \'_\'`
		}
	}
	email := strings.ToLower(body.Email)
	if len(email) == 0 {
		response.ErrFor["email"] = "required"
	} else {
		r, err := regexp.MatchString(`^[a-zA-Z0-9\-\_\.\+]+@[a-zA-Z0-9\-\_\.]+\.[a-zA-Z0-9\-\_]+$`, email)
		if err != nil {
			println(err.Error())
		}
		if !r {
			response.ErrFor["email"] = `invalid email format`
		}
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail(c)
		return
	}


	if len(user.IsActive) == 0 {
		user.IsActive = "no"
	}

	// duplicateUsernameCheck
	// duplicateEmailCheck
	{
		us := User{} // todo pool
		err = collection.Find(bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}).One(&us)
		if err != nil {
			response.Errors = append(response.Errors, "username or email already exist")
			response.Fail(c)
			return
		}
	}

	user.Username = body.Username
	user.Email = body.Email

	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}
	// TODO  patch admin and account
	response.Success = true
	c.JSON(http.StatusOK, response)





	response := Response{} // todo sync.Pool
	response.Errors = []string{}
	response.ErrFor = make(map[string]string)
	defer response.Recover(c)
	var body struct {
		Confirm   string  `json:"confirm"`
		Password string  `json:"newPassword"`
	}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&body)

	// validate
	if len(body.Password) == 0 {
		response.ErrFor["newPassword"] = "required"
	}
	if len(body.Confirm) == 0 {
		response.ErrFor["confirm"] = "required"
	} else if body.Password != body.Confirm {
		response.Errors = append(response.Errors, "Passwords do not match.")
	}

	if response.HasErrors() {
		response.Fail(c)
		return
	}

	// patchUser
	db := getMongoDBInstance()
	defer db.Session.Close()
	collection := db.C(USERS)
	user := User{}
	err = collection.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			panic(err)
		}
		response.Errors = append(response.Errors, "User wasn't found.")
		response.Fail(c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Errors = append(response.Errors, err.Error()) // TODO don't like that this error goes to client
		response.Fail(c)
		return
	}

	user.Password = string(hashedPassword)
	err = collection.UpdateId(user.ID, user)
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		response.Fail(c)
		return
	}

	response.Success = true
	c.JSON(http.StatusOK, response)
}
*/

