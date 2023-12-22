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

type PathMapping struct {
	function   HandlerFunction
	pathParams map[string]int
	method     string
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
	if strings.Contains(path, ":") {
		var pathName string = ""
		for index, value := range strings.Split(path, "/") {
			if value == "" {
				continue
			}
			if value[0] == ':' {
				pathParam[value[1:]] = index - 1
			} else {
				pathName += "/" + value
			}
		}
		path = pathName
	}
	if path == "" {
		path = "/"
	}
	pathMapping.function = function
	pathMapping.pathParams = pathParam
	r.routeMapping[path] = append(r.routeMapping[path], *pathMapping)
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func valueIn(i int, dict map[string]int) bool {
	for _, value := range dict {
		if i == value {
			return true
		}
	}
	return false
}

func extractPath(urlArray []string, param PathMapping) string {
	var path string
	for i, value := range urlArray {
		if valueIn(i, param.pathParams) {
			continue
		}
		path += "/" + value
	}
	return path
}

func getSignleSlashFunction(params []PathMapping, urlArray []string) HandlerFunction {
	var function HandlerFunction

	for _, param := range params {
		if len(urlArray) == len(param.pathParams) {
			return param.function
		}
	}
	return function
}

func removeBlankStrings(array []string) []string {
	var copyArray []string
	for _, value := range array {
		if value != "" {
			copyArray = append(copyArray, value)
		}
	}
	return copyArray

}

func mapKey(m map[string]int, value int) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func getParams(urlArray []string, param PathMapping) map[string]string {

	paramValueMap := make(map[string]string)
	for i, value := range urlArray {
		if valueIn(i, param.pathParams) {
			key, ok := mapKey(param.pathParams, i)
			if ok {
				paramValueMap[key] = value
			}
		}
	}
	return paramValueMap
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	pathArray := strings.Split(url, "/")
	pathArray = removeBlankStrings(pathArray)

	var function HandlerFunction

	index := 0
	var slashMap []PathMapping

	for path, avialableParams := range r.routeMapping {
		if function != nil {
			break
		}
		if path == "/" {
			slashMap = avialableParams
		}

		for _, param := range avialableParams {
			finalPath := extractPath(pathArray, param)
			paramValueMap := getParams(pathArray, param)
			if finalPath == path && (len(paramValueMap) == len(param.pathParams)) {
				fmt.Println(finalPath, paramValueMap, param.pathParams)
				function = param.function
				break
			}
		}
		if index == len(r.routeMapping)-1 && len(slashMap) > 0 {
			fmt.Println("here")
			function = getSignleSlashFunction(slashMap, pathArray)
		}
		index++
	}
	if function != nil {
		function(w, req)
	} else {
		default404Response(w, req)
	}
}
