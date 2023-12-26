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

func HomeSomePathPatch(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/home/:somePath, patch")
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

func HomePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home post")
}

func main() {
	r := CreateRouter()
	r.addRoute("/home", Home, GET)
	r.addRoute("/home", HomePost, POST)
	r.addRoute("/home", Random, POST)
	r.addRoute("/home", HomePost, POST)
	r.addRoute("/", Index, GET)
	r.addRoute("/:anyParam", Hello, POST)
	r.addRoute("/:anyParam", Hello, GET)
	r.addRoute("/home/:somePath", HomeSomePath, POST)
	r.addRoute("/home/:somePath", HomeSomePathPatch, PATCH)
	r.addRoute("/random", Random, PATCH)
	r.addRoute("/random/:params", RandomWithParams, DELETE)
	r.addRoute("/random/:params", RandomWithParams, POST)
	r.addRoute("/:somePath/asdf/:nothing", HomeSomePath2, CONNECT)
	r.addRoute("/asdf/:fasdfasd/:jflakdsfas", HomeSomePath2, CONNECT)
	r.addRoute("/file/:fileName/:anotherParam", FileNotRandom, OPTIONS)
	r.addRoute("/file/:fileName/random/:somethingElse", FileRandom, TRACE)

	fmt.Println("Started listening on 0.0.0.0:8080")
	for routeHandler, function := range r.routeMapping {
		fmt.Printf("%s, %v, %s, %s\n",
			routeHandler.route, routeHandler.pathParams, routeHandler.http_method, GetFunctionName(function))
	}
	fmt.Println("------------------------------->")
	err := http.ListenAndServe("0.0.0.0:8080", r)

	if err != nil {
		fmt.Printf("failed to start server: %s \n", err)
	}
}
