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
					pathParam[value[1:]] = slashIndex
				} else {
					pathName += "/" + value
				}
				slashIndex++
			}
			path = pathName
		}
	}
	if path[len(path)-1] == '/' && len(path) != 1 {
		path = path[:len(path)-1]
	}
	pathMapping.function = function
	pathMapping.pathParams = pathParam
	r.routeMapping[path] = append(r.routeMapping[path], *pathMapping)

}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func ValueIn(i int, dict map[string]int) bool {
	for _, value := range dict {
		if i == value {
			return true
		}
	}
	return false
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	urlArray := strings.Split(url, "/")
	urlArray = urlArray[1:]

	var function HandlerFunction
	var found bool
	var finalPath string

	for path, avialableParams := range r.routeMapping {
		if found {
			break
		}
		if path == url {
			function = avialableParams[0].function
			break
		}

		for _, paramPostion := range avialableParams {
			if found {
				break
			}
			finalPath = ""
			for i, value := range urlArray {
				if ValueIn(i, paramPostion.pathParams) {
					continue
				}
				finalPath += "/" + value
			}
			fmt.Printf("final path %s | %s \n", finalPath, path)
			if strings.Compare(finalPath, path) == 0 {
				function = paramPostion.function
			}
		}

	}
	if function != nil {
		function(w, req)
	}

}

func ParamVars(req *http.Request) map[string]string {
	path := req.URL.String()
	params := make(map[string]string)
	fmt.Println(path)
	return params
}
