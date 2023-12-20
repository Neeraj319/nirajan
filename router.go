package main

import (
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

type PathMapping struct {
	function   HandlerFunction
	pathParams map[string]int
}

func CreatePathMapping() *PathMapping {
	return &PathMapping{pathParams: make(map[string]int)}
}

type SimpleRouter struct {
	routeMapping map[string][]PathMapping
	Response404  HandlerFunction
}

func CreateRouter() *SimpleRouter {
	router := SimpleRouter{
		routeMapping: make(map[string][]PathMapping),
		Response404:  default404Response,
	}
	return &router
}

func (r *SimpleRouter) addRoute(path string, function HandlerFunction) {
	pathParam := make(map[string]int)
	pathMapping := CreatePathMapping()
	var slashIndex int
	if strings.Contains(path, ":") {
		pathParams := strings.Split(path, ":")
		if len(pathParams) == 2 {
			pathParam[pathParams[1]] = slashIndex
			path = pathParams[0]
		} else {
			var pathName string = ""
			for _, value := range strings.Split(path, "/") {
				if value == "" {
					continue
				}
				if value[0] == ':' {
					// pathParamsNames = append(pathParamsNames, value[1:])
					pathParam[value[1:]] = slashIndex
					slashIndex++
				} else {
					pathName += "/" + value
				}
			}
			path = pathName
		}
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	pathMapping.function = function
	pathMapping.pathParams = pathParam
	r.routeMapping[path] = append(r.routeMapping[path], *pathMapping)
	fmt.Println(r.routeMapping)

}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var function HandlerFunction
	path := req.URL.String()
	mappings := r.routeMapping[path]
	fmt.Println(path)

	if len(mappings) == 0 && r.Response404 != nil {
		default404Response(w, req)
		return
	}
	fmt.Println(mappings)
	if len(mappings) == 1 {
		function = mappings[0].function
		function(w, req)
	}
}

func ParamVars(req *http.Request) map[string]string {
	path := req.URL.String()
	params := make(map[string]string)
	fmt.Println(path)
	return params
}
