package main

import (
	// "fmt"
	"fmt"
	"github.com/gorilla/schema"
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

var decoder = schema.NewDecoder()

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

type HandlerFunctionPathParams func(w http.ResponseWriter, r *http.Request, paramsStruct interface{})
type HandlerFunctionNoPathParams func(w http.ResponseWriter, r *http.Request)

func default404Resp(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
}

func defaultMethodNotAllowedResp(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(405)
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
	routeMapping         map[*RouteHandler]interface{}
	NotFoundResp         HandlerFunctionNoPathParams
	MethodNotAllowedResp HandlerFunctionNoPathParams
}

func CreateRouter() *SimpleRouter {
	router := SimpleRouter{
		routeMapping:         make(map[*RouteHandler]interface{}),
		NotFoundResp:         default404Resp,
		MethodNotAllowedResp: defaultMethodNotAllowedResp,
	}
	return &router
}

func removeSimilarRoute(r *map[*RouteHandler]interface{}, routeHandler RouteHandler) {
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

func valueIn(i int, dict map[string]int) bool {
	for _, value := range dict {
		if i == value {
			return true
		}
	}
	return false
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

func validateHanlder(function interface{}) {
	v := reflect.TypeOf(function)
	if v.Kind() != reflect.Func {
		panic("param: function (must be of function type)")
	}

	validateHandlerParams(function)
}

func validateHandlerParams(function interface{}) {
	functionName := GetFunctionName(function)

	v := reflect.TypeOf(function)

	_, file, line, _ := runtime.Caller(1)

	errorTemplate := "[%s:%d] %s: %s"

	if !(v.NumIn() == 3 || v.NumIn() == 2) {
		panic(fmt.Sprintf(errorTemplate, file, line, functionName, "handler function params must be 2 or 3"))
	}

	if v.In(0).Kind() != reflect.Interface || v.In(0).String() != "http.ResponseWriter" {
		panic(fmt.Sprintf(errorTemplate, file, line, functionName, "handler first argument must be of type `http.ResponseWriter`"))
	}

	if v.In(1).Kind() != reflect.Ptr || v.In(1).Elem().String() != "http.Request" {
		panic(fmt.Sprintf(errorTemplate, file, line, functionName, "handler second argument must be a pointer to `http.Request`"))
	}

	if v.NumIn() == 3 && v.In(2).Kind() != reflect.Struct {
		panic(fmt.Sprintf(errorTemplate, file, line, functionName, "handler third argument must be a struct"))
	}
}

func validateHandlerStruct() {

}

func (r *SimpleRouter) addRoute(path string, function interface{}, http_method HTTP_METHOD) {
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

	validateHanlder(function)

	routeHandler.pathParams = pathParams
	routeHandler.route = path
	removeSimilarRoute(&r.routeMapping, *routeHandler)
	r.routeMapping[routeHandler] = function
}

func GetFunctionName(temp interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()
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

func getSingleSlashHandlers(params []*RouteHandler, urlArray []string) []*RouteHandler {

	routeHandlers := make([]*RouteHandler, 0)
	for _, routeHandlerObj := range params {
		if len(urlArray) == len(routeHandlerObj.pathParams) {
			routeHandlers = append(routeHandlers, routeHandlerObj)
		}
	}
	return routeHandlers
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

func getParams(urlArray []string, paramIndex map[string]int) map[string][]string {
	paramValueMap := make(map[string][]string)
	for i, value := range urlArray {
		if valueIn(i, paramIndex) {
			key, ok := mapKey(paramIndex, i)
			if ok {
				paramValueMap[key] = append(paramValueMap[key], value)
			}
		}
	}
	return paramValueMap
}

func getPossibleRouteHandlers(routeMapping map[*RouteHandler]interface{}, pathArray []string) []*RouteHandler {
	var slashMap []*RouteHandler
	var possibleRouteHandlers []*RouteHandler
	for routeObj := range routeMapping {
		if routeObj.route == "/" {
			slashMap = append(slashMap, routeObj)
			continue
		}

		finalPath := extractPath(pathArray, routeObj)
		paramValueMap := getParams(pathArray, routeObj.pathParams)
		if finalPath == routeObj.route && (len(paramValueMap) == len(routeObj.pathParams)) {
			possibleRouteHandlers = append(possibleRouteHandlers, routeObj)
		}
	}

	/*
		if there are routes with only "/" and if route is not found then check in "/" routes
		for this only check for the length of the param
	*/

	if len(possibleRouteHandlers) == 0 && len(slashMap) > 0 {
		possibleRouteHandlers = getSingleSlashHandlers(slashMap, pathArray)
	}
	return possibleRouteHandlers
}

func callHandler(handler interface{}, pathArray []string, handerObj *RouteHandler, defaultArgs []reflect.Value) {
	typeFunction := reflect.TypeOf(handler)
	function := reflect.ValueOf(handler)
	switch typeFunction.NumIn() {
	case 2:
		function.Call(defaultArgs)
		return
	case 3:
		params := getParams(pathArray, handerObj.pathParams)

		v := reflect.New(typeFunction.In(2))
		decoder.Decode(v.Interface(), params)

		args := append(defaultArgs, v.Elem())
		function.Call(args)
		return
	}
}

func (r *SimpleRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	pathArray := strings.Split(url, "/")
	pathArray = removeBlankStrings(pathArray)

	possibleRouteHandlers := getPossibleRouteHandlers(r.routeMapping, pathArray)

	if len(possibleRouteHandlers) == 0 {
		r.NotFoundResp(w, req)
		return
	}
	for _, routeObj := range possibleRouteHandlers {
		if routeObj.http_method.String() == req.Method {
			function := r.routeMapping[routeObj]
			args := []reflect.Value{
				reflect.ValueOf(w),
				reflect.ValueOf(req),
			}
			callHandler(function, pathArray, routeObj, args)
			return
		}
	}
	r.MethodNotAllowedResp(w, req)
}
