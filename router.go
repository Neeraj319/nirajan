package main

import (
	// "fmt"
	"fmt"
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
	pathParams   map[string][]int
}

func (r *SimpleRouter) addRoute(path string, function HandlerFunction) {
	if strings.Contains(path, ":") {
		pathParams := strings.Split(path, ":")
		if len(pathParams) == 2 {
			params := make([]int, 1)
			params[0] = 1
			r.pathParams[pathParams[0]] = params
		} else {
			var params []int
			var paramsPathName string
			for i, value := range path {
				if rune(value) == ':' {
					params = append(params, i-1)
				}
			}
			for _, value := range strings.Split(path, "/") {
				if (value) != "" && value[0] != ':' {
					paramsPathName += value
				}
			}
			r.pathParams[paramsPathName] = params
			fmt.Println(r.pathParams)
		}
	}
	r.routeMapping[path] = function
}

func CreateRouter() *SimpleRouter {
	return &SimpleRouter{routeMapping: make(map[string]HandlerFunction), Response404: default404Response, pathParams: make(map[string][]int)}
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	fmt.Println(url)
	function := r.routeMapping[url]

	if function == nil && r.Response404 != nil {
		default404Response(w, req)
		return
	}
	function(w, req)
}
