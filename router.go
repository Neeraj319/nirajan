package main

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type HandlerFunction func(w http.ResponseWriter, r *http.Request)

func default404Response(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
}

type SimpleRouter struct {
	routeMapping map[string]HandlerFunction
	Response404  HandlerFunction
}

func (r *SimpleRouter) addRoute(path string, function HandlerFunction) {
	r.routeMapping[path] = function
}

func CreateRouter() *SimpleRouter {
	return &SimpleRouter{routeMapping: make(map[string]HandlerFunction), Response404: default404Response}
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	function := r.routeMapping[url]

	if function == nil && r.Response404 != nil {
		default404Response(w, req)
		return
	}
	function(w, req)
}
