package main

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type HTTP_METHOD int

const (
	POST HTTP_METHOD = iota
	GET
	HEAD
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	PATCH
)

func (me HTTP_METHOD) String() string {
	methodNames := [...]string{
		"POST",
		"GET",
		"HEAD",
		"PUT",
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
		"PATCH",
	}

	if int(me) < len(methodNames) {
		return methodNames[me]
	}
	return ""
}

type HandlerFunction func(w http.ResponseWriter, r *http.Request)

func default404Resp(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
}

func defaultMethodNotAllowedResp(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(405)
}

type RouteHandler struct {
	function   HandlerFunction
	pathParams map[string]int
	method     HTTP_METHOD
}

func createRouteHandler(method HTTP_METHOD) *RouteHandler {
	return &RouteHandler{pathParams: make(map[string]int), method: method}
}

type SimpleRouter struct {
	routeMapping         map[string][]RouteHandler
	NotFoundResp         HandlerFunction
	MethodNotAllowedResp HandlerFunction
}

func CreateRouter() *SimpleRouter {
	router := SimpleRouter{
		routeMapping:         make(map[string][]RouteHandler),
		NotFoundResp:         default404Resp,
		MethodNotAllowedResp: defaultMethodNotAllowedResp,
	}
	return &router
}

func (r *SimpleRouter) addRoute(path string, function HandlerFunction, method HTTP_METHOD) {
	pathParams := make(map[string]int)
	routeHandler := createRouteHandler(method)
	if strings.Contains(path, ":") {
		var pathName string = ""
		for index, value := range strings.Split(path, "/") {
			if value == "" {
				continue
			}
			if value[0] == ':' {
				pathParams[value[1:]] = index - 1
			} else {
				pathName += "/" + value
			}
		}
		path = pathName
	}
	if path == "" {
		path = "/"
	}
	routeHandler.function = function
	routeHandler.pathParams = pathParams
	r.routeMapping[path] = append(r.routeMapping[path], *routeHandler)
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

func extractPath(urlArray []string, param RouteHandler) string {
	var path string
	for i, value := range urlArray {
		if valueIn(i, param.pathParams) {
			continue
		}
		path += "/" + value
	}
	return path
}

func getSingleSlashHandler(params []RouteHandler, urlArray []string) RouteHandler {

	for _, routeHandlerObj := range params {
		if len(urlArray) == len(routeHandlerObj.pathParams) {
			return routeHandlerObj
		}
	}
	return RouteHandler{}
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

func getParams(urlArray []string, param RouteHandler) map[string]string {

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

	index := 0
	var slashMap []RouteHandler
	var routeHandler RouteHandler

	for pathString, routeHandlers := range r.routeMapping {
		if routeHandler.function != nil {
			break
		}
		if pathString == "/" {
			slashMap = routeHandlers
		}

		for _, routeObj := range routeHandlers {
			finalPath := extractPath(pathArray, routeObj)
			paramValueMap := getParams(pathArray, routeObj)
			if finalPath == pathString && (len(paramValueMap) == len(routeObj.pathParams)) {
				routeHandler = routeObj
				break
			}
		}
		if index == len(r.routeMapping)-1 && len(slashMap) > 0 {
			routeHandler = getSingleSlashHandler(slashMap, pathArray)
		}
		index++
	}
	if routeHandler.function != nil {
		if routeHandler.method.String() != req.Method {
			r.MethodNotAllowedResp(w, req)
		}
		routeHandler.function(w, req)
	} else {
		r.NotFoundResp(w, req)
	}
}
