# What it is?

**Gouter** is implementation of http router created for understanding principle of trie data structure and prefix find algorithms. It has a pretty simple code to understand, but the functionality is poor. The idea of this project is the most readable code for understanding basic functionality of Radix tree.

# Usage

```go
package main

import (
	"github.com/furdarius/gouter"
	"io"
	"net/http"
)

func Index(w http.ResponseWriter, req *http.Request, params gouter.Params) {
	io.WriteString(w, "Index Method")
}

func Show(w http.ResponseWriter, req *http.Request, params gouter.Params) {
	// gouter.Params is just map where key is param name
	io.WriteString(w, "Show Method. ID: "+params["id"])
}

func Edit(w http.ResponseWriter, req *http.Request, params gouter.Params) {
	io.WriteString(w, "Edit Method")
}

func main() {
	gouterInstance := gouter.New()

	gouterInstance.Add("GET", "/", func(w http.ResponseWriter, req *http.Request, params gouter.Params) {
		io.WriteString(w, "My url is "+req.URL.String())
	})

	gouterInstance.Add("GET", "/list", Index)

	gouterInstance.Add("GET", "/list/:id", Show)

	gouterInstance.Add("GET", "/list/:id/edit", Edit)

	// You can set own 404 error handler
	gouterInstance.SetNotFoundHandler(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "404 custom error!", http.StatusNotFound)
	})

	http.ListenAndServe(":8080", gouterInstance)
}

```