package main

import (
	// "fmt"
	"fmt"
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

func valueIn(i int, dict map[string]int) bool {
	for _, value := range dict {
		if i == value {
			return true
		}
	}
	return false
}

type RouteHandler struct {
	pathParams  map[string]int
	route       string
	http_method HTTP_METHOD
}

func createRouteHandler(http_method HTTP_METHOD) *RouteHandler {
	return &RouteHandler{pathParams: make(map[string]int), http_method: http_method}
}

type SimpleRouter struct {
	routeMapping         map[*RouteHandler]HandlerFunction
	NotFoundResp         HandlerFunction
	MethodNotAllowedResp HandlerFunction
}

func CreateRouter() *SimpleRouter {
	router := SimpleRouter{
		routeMapping:         make(map[*RouteHandler]HandlerFunction),
		NotFoundResp:         default404Resp,
		MethodNotAllowedResp: defaultMethodNotAllowedResp,
	}
	return &router
}

func removeSimilarRoute(r *map[*RouteHandler]HandlerFunction, routeHandler RouteHandler) {
	actualMap := *r
	keysToDelete := make([]*RouteHandler, 0)

	for routeObj := range actualMap {
		if isRouteSimilar(routeObj, routeHandler) {
			keysToDelete = append(keysToDelete, routeObj)
		}
	}

	for _, key := range keysToDelete {
		delete(actualMap, key)
	}
}

func isRouteSimilar(routeObj *RouteHandler, routeHandler RouteHandler) bool {
	return routeObj.route == routeHandler.route &&
		routeObj.http_method == routeHandler.http_method &&
		len(routeObj.pathParams) == len(routeHandler.pathParams) &&
		areAllParamsPresent(routeObj.pathParams, routeHandler.pathParams)
}

func areAllParamsPresent(objParams map[string]int, handlerParams map[string]int) bool {
	for _, param := range objParams {
		if !valueIn(param, handlerParams) {
			return false
		}
	}
	return true
}

func (r *SimpleRouter) addRoute(path string, function HandlerFunction, http_method HTTP_METHOD) {
	pathParams := make(map[string]int)
	routeHandler := createRouteHandler(http_method)
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
	routeHandler.pathParams = pathParams
	routeHandler.route = path
	removeSimilarRoute(&r.routeMapping, *routeHandler)
	r.routeMapping[routeHandler] = function
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func extractPath(urlArray []string, param *RouteHandler) string {
	var path string
	for i, value := range urlArray {
		if valueIn(i, param.pathParams) {
			continue
		}
		path += "/" + value
	}
	return path
}

func getSingleSlashHandler(params []*RouteHandler, urlArray []string) (*RouteHandler, bool) {

	for _, routeHandlerObj := range params {
		if len(urlArray) == len(routeHandlerObj.pathParams) {
			return routeHandlerObj, true
		}
	}
	return &RouteHandler{}, false
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

func getParams(urlArray []string, param *RouteHandler) map[string]string {

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

	var slashMap []*RouteHandler
	var possibleRouteHandlers []*RouteHandler

	for routeObj := range r.routeMapping {
		if routeObj.route == "/" {
			slashMap = append(slashMap, routeObj)
			continue
		}

		finalPath := extractPath(pathArray, routeObj)
		paramValueMap := getParams(pathArray, routeObj)
		if finalPath == routeObj.route && (len(paramValueMap) == len(routeObj.pathParams)) {
			possibleRouteHandlers = append(possibleRouteHandlers, routeObj)
		}
	}
	/*
		if there are routes with only "/" and if route is not found then check in "/" routes
		for this only check for the length of the param
	*/

	if len(possibleRouteHandlers) == 0 && len(slashMap) > 0 {
		routeObj, ok := getSingleSlashHandler(slashMap, pathArray)
		if ok {
			possibleRouteHandlers = append(possibleRouteHandlers, routeObj)
		}
	}
	if len(possibleRouteHandlers) == 0 {
		r.NotFoundResp(w, req)
		return
	}
	for _, routeObj := range possibleRouteHandlers {
		fmt.Println("possible", routeObj.http_method, routeObj.route, routeObj.pathParams)
		if routeObj.http_method.String() == req.Method {
			r.routeMapping[routeObj](w, req)
			return
		}
	}
	r.MethodNotAllowedResp(w, req)
}
