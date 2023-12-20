package main

import (
	"fmt"
	"net/http"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	parmVars := ParamVars(r)
	fmt.Println(parmVars)
	fmt.Fprintf(w, "hello")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello ")
}

func main() {
	r := CreateRouter()
	r.addRoute("/home", index)
	r.addRoute("/", index)
	r.addRoute("/:hello", index)
	r.addRoute("/home/:somePath", fileHandler)
	r.addRoute("/file/:fileName/:somethingElse", fileHandler)
	r.addRoute("/file/:fileName/random/:somethingElse", fileHandler)

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Printf("failed to start server: %s \n", err)
	}
}
