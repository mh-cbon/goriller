# goriller

[![travis Status](https://travis-ci.org/mh-cbon/goriller.svg?branch=master)](https://travis-ci.org/mh-cbon/goriller) [![Appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/goriller?branch=master&svg=true)](https://ci.appveyor.com/projects/mh-cbon/goriller) [![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/goriller)](https://goreportcard.com/report/github.com/mh-cbon/goriller) [![GoDoc](https://godoc.org/github.com/mh-cbon/goriller?status.svg)](http://godoc.org/github.com/mh-cbon/goriller) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Package goriller generate gorilla routers.


# TOC
- [Install](#install)
  - [Usage](#usage)
    - [$ goriller -help](#-goriller--help)
  - [Cli examples](#cli-examples)
- [API example](#api-example)
  - [Anootations](#anootations)
  - [> demo/main.go](#-demomaingo)
  - [> demo/controllergoriller.go](#-democontrollergorillergo)
  - [> demo/controllergorillerrpc.go](#-democontrollergorillerrpcgo)
- [Recipes](#recipes)
  - [Release the project](#release-the-project)
- [History](#history)

# Install
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/goriller
cd $GOPATH/src/github.com/mh-cbon/goriller
git clone https://github.com/mh-cbon/goriller.git .
glide install
go install
```

## Usage

#### $ goriller -help
```sh
goriller 0.0.0

Usage
	goriller [-p name] [-mode name] [...types]

  types:  A list of types such as src:dst.
          A type is defined by its package path and its type name,
          [pkgpath/]name
          If the Package path is empty, it is set to the package name being generated.
          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName 
  -p:     The name of the package output.
  -mode:  The generation mode std|rpc.
```

## Cli examples

```sh
# Create a goriller binder version of JSONTomates to HTTPTomates
goriller *JSONTomates:HTTPTomates
# Create a goriller binder version of JSONTomates to HTTPTomates to stdout
goriller -p main - JSONTomates:HTTPTomates
```

# API example

Following example demonstates a program using it to generate a goriller binder of a type.

#### Anootations

`goriller` reads and interprets annotations on `struct` and `methods`.

The `struct` annotations are used as default for the `methods` annotations.

| Name | Description |
| --- | --- |
| @route | The route path such as `/{param}` |
| @name | The route name `name` |
| @host | The route name `host` |
| @methods | The route methods `GET,POST,PUT` |
| @schemes | The route methods `http, https` |

#### > demo/main.go
```go
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
```

Following code is the generated implementation of the goriller binder.

#### > demo/controllergoriller.go
```go
package main

// file generated by
// github.com/mh-cbon/goriller
// do not edit

import (
	"github.com/gorilla/mux"
)

// ControllerGoriller is a goriller of *ControllerHTTPGen.
// ControllerHTTPGen is an httper of *ControllerJSONGen.
// ControllerJSONGen is jsoner of *Controller.
// Controller of some resources.
type ControllerGoriller struct {
	embed *ControllerHTTPGen
}

// NewControllerGoriller constructs a goriller of *ControllerHTTPGen
func NewControllerGoriller(embed *ControllerHTTPGen) *ControllerGoriller {
	ret := &ControllerGoriller{
		embed: embed,
	}
	return ret
}

// Bind the given router.
func (t ControllerGoriller) Bind(router *mux.Router) {
	router.HandleFunc("/{id}", t.embed.GetByID).Methods("GET")
	router.HandleFunc("/{id}", t.embed.UpdateByID).Methods("PUT", "POST")
	router.HandleFunc("/{id}", t.embed.DeleteByID).Methods("DELETE")

}
```

Following code is the generated implementation of the goriller binder in an rpc fashion.

#### > demo/controllergorillerrpc.go
```go
package main

// file generated by
// github.com/mh-cbon/goriller
// do not edit

import (
	"github.com/gorilla/mux"
)

// ControllerGorillerRPC is a goriller of *ControllerHTTPGen.
// ControllerHTTPGen is an httper of *ControllerJSONGen.
// ControllerJSONGen is jsoner of *Controller.
// Controller of some resources.
type ControllerGorillerRPC struct {
	embed *ControllerHTTPGen
}

// NewControllerGorillerRPC constructs a goriller of *ControllerHTTPGen
func NewControllerGorillerRPC(embed *ControllerHTTPGen) *ControllerGorillerRPC {
	ret := &ControllerGorillerRPC{
		embed: embed,
	}
	return ret
}

// Bind the given router.
func (t ControllerGorillerRPC) Bind(router *mux.Router) {
	router.HandleFunc("GetByID", t.embed.GetByID)
	router.HandleFunc("UpdateByID", t.embed.UpdateByID)
	router.HandleFunc("DeleteByID", t.embed.DeleteByID)
	router.HandleFunc("TestVars1", t.embed.TestVars1)
	router.HandleFunc("TestCookier", t.embed.TestCookier)
	router.HandleFunc("TestSessionner", t.embed.TestSessionner)
	router.HandleFunc("TestRPCer", t.embed.TestRPCer)

}
```


# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
