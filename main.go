package main

import (
	"fmt"
	"net/http"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	parmVars := ParamVars(r)
	fmt.Println(parmVars)
	fmt.Fprintf(w, "file handler no random")
}

func fileHandler2(w http.ResponseWriter, r *http.Request) {
	parmVars := ParamVars(r)
	fmt.Println(parmVars)
	fmt.Fprintf(w, "file handler with random")
}
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is index")
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is home ")
}
func homeSomePath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is with some path ")
}
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello ")
}
func main() {
	r := CreateRouter()
	r.addRoute("/home", home)
	r.addRoute("/", index)
	r.addRoute("/:hello", hello)
	r.addRoute("/home/:somePath", homeSomePath)
	r.addRoute("/file/:fileName/:somethingElse", fileHandler)
	r.addRoute("/file/:fileName/random/:somethingElse", fileHandler2)

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Printf("failed to start server: %s \n", err)
	}
}
