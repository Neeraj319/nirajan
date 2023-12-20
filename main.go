package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["file_name"]
	fmt.Fprintf(w, "File Name: %s", fileName)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello ")
}

func main() {
	r := CreateRouter()
	r.addRoute("/", index)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("faile to start server: %s \n", err)
	}

}
