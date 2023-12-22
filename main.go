package main

import (
	"fmt"
	"net/http"
)

func fileNotRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/file/:fileName/:anotherParam")
}

func fileRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/file/:fileName/random/:somethingElse")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/")
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/home")
}

func homeSomePath2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/:somePath/asdf/:nothing")
}

func homeSomePath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/home/:somePath")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/:anyParam")
}

func random(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/random")
}

func randomWithParams(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/random/:params")
}

func main() {
	r := CreateRouter()
	r.addRoute("/home", home)
	r.addRoute("/", index)
	r.addRoute("/:anyParam", hello)
	r.addRoute("/home/:somePath", homeSomePath)
	r.addRoute("/random", random)
	r.addRoute("/random/:params", randomWithParams)
	r.addRoute("/:somePath/asdf/:nothing", homeSomePath2)
	r.addRoute("/file/:fileName/:anotherParam", fileNotRandom)
	r.addRoute("/file/:fileName/random/:somethingElse", fileRandom)

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Printf("failed to start server: %s \n", err)
	}
}
