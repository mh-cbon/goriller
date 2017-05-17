package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	httper "github.com/mh-cbon/httper/lib"
)

//go:generate lister *Tomate:TomatesGen
//go:generate channeler TomatesGen:TomatesSyncGen

//go:generate jsoner -mode gorilla *Controller:ControllerJSONGen
//go:generate httper -mode gorilla *ControllerJSONGen:ControllerHTTPGen
//go:generate goriller *ControllerHTTPGen:ControllerGoriller
//go:generate goriller -mode rpc *ControllerHTTPGen:ControllerGorillerRPC

func main() {

	backend := NewTomatesSyncGen()
	backend.Push(&Tomate{Name: "Red"})

	jsoner := NewControllerJSONGen(NewController(backend), nil)
	httper := NewControllerHTTPGen(jsoner, nil)

	router := mux.NewRouter()
	NewControllerGoriller(httper).Bind(router)
	http.Handle("/", router)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	time.Sleep(1 * time.Millisecond)

	req, err := http.Get("http://localhost:8080/0")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	io.Copy(os.Stdout, req.Body)
}

// Tomate is about red vegetables to make famous italian food.
type Tomate struct {
	ID   int
	Name string
}

// GetID return the ID of the Tomate.
func (t *Tomate) GetID() int {
	return t.ID
}

// Controller of some resources.
type Controller struct {
	backend *TomatesSyncGen
}

// NewController ...
func NewController(backend *TomatesSyncGen) *Controller {
	return &Controller{
		backend: backend,
	}
}

// GetByID ...
// @route /{id}
// @methods GET
func (t *Controller) GetByID(urlID int) *Tomate {
	return t.backend.Filter(FilterTomatesGen.ByID(urlID)).First()
}

// UpdateByID ...
// @route /{id}
// @methods PUT,POST
func (t *Controller) UpdateByID(urlID int, reqBody *Tomate) *Tomate {
	var ret *Tomate
	t.backend.Filter(func(v *Tomate) bool {
		if v.ID == urlID {
			v.Name = reqBody.Name
			ret = v
		}
		return true
	})
	return ret
}

// DeleteByID ...
// @route /{id}
// @methods DELETE
func (t *Controller) DeleteByID(REQid int) bool {
	return t.backend.Remove(&Tomate{ID: REQid})
}

// TestVars1 ...
func (t *Controller) TestVars1(w http.ResponseWriter, r *http.Request) {
}

// TestCookier ...
func (t *Controller) TestCookier(c httper.Cookier) {
}

// TestSessionner ...
func (t *Controller) TestSessionner(s httper.Sessionner) {
}

// TestRPCer ...
func (t *Controller) TestRPCer(id int) bool {
	return false
}
