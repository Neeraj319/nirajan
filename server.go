package main

import (
	"fmt"
	"net/http"
)

func FileNotRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/file/:fileName/:anotherParam")
}

func FileRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/file/:fileName/random/:somethingElse")
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/")
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/home")
}

func HomeSomePath2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/:somePath/asdf/:nothing")
}

func HomeSomePath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/home/:somePath")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/:anyParam")
}

func Random(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/random")
}

func RandomWithParams(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/random/:params")
}

func main() {
	r := CreateRouter()
	r.addRoute("/home", Home)
	r.addRoute("/", Index)
	r.addRoute("/:anyParam", Hello)
	r.addRoute("/home/:somePath", HomeSomePath)
	r.addRoute("/random", Random)
	r.addRoute("/random/:params", RandomWithParams)
	r.addRoute("/:somePath/asdf/:nothing", HomeSomePath2)
	r.addRoute("/file/:fileName/:anotherParam", FileNotRandom)
	r.addRoute("/file/:fileName/random/:somethingElse", FileRandom)

	fmt.Println("Started listening on 0.0.0.0:8080")
	err := http.ListenAndServe("0.0.0.0:8080", r)

	if err != nil {
		fmt.Printf("failed to start server: %s \n", err)
	}
}
