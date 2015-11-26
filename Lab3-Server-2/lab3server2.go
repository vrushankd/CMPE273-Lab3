package main

import (
	// Standard library packages
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	// Third party packages
	"github.com/julienschmidt/httprouter"
)

//UserController structure
type UserController struct{}

//NewUserController function
func NewUserController() *UserController {
	return &UserController{}
}

//Data struct
type Data struct {
	Datakey   int    `json:"key"`
	Datavalue string `json:"value"`
}

var datamap = make(map[int]string)

func main() {
	// Instantiate a new router
	router := httprouter.New()

	// Get a controller instance
	controller := NewUserController()

	// Add handlers
	router.GET("/keys/:id", controller.keybyID)
	router.GET("/keys", controller.keyAll)
	router.PUT("/keys/:id/:value", controller.putkey)

	// Expose the server at port 3000
	http.ListenAndServe(":3001", router)
}

func (uc UserController) keybyID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	strtointID, _ := strconv.Atoi(id)

	var (
		k int
		v string
	)

	for k, v = range datamap {
		if k == strtointID {
			v = datamap[k]
			break
		}
	}
	uj := Data{
		Datakey:   k,
		Datavalue: v,
	}

	result, _ := json.Marshal(uj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", result)
}

func (uc UserController) keyAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	type Users []Data
	var useraarr = make(Users, 0)
	//var uj Data
	for k, v := range datamap {
		uj := Data{k, v}
		useraarr = append(useraarr, uj)
	}

	result, _ := json.Marshal(useraarr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", result)

}

func (uc UserController) putkey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	val := p.ByName("value")
	strtoIntID, _ := strconv.Atoi(id)

	datamap[strtoIntID] = val
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
