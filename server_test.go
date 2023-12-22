package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRoute(t *testing.T, route string, function HandlerFunction) {
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(function)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != route {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), route)
	}
}

func TestAllRoutes(t *testing.T) {
	routes := map[string]HandlerFunction{
		"/home":                                 Home,
		"/":                                     Index,
		"/:anyParam":                            Hello,
		"/home/:somePath":                       HomeSomePath,
		"/random":                               Random,
		"/random/:params":                       RandomWithParams,
		"/:somePath/asdf/:nothing":              HomeSomePath2,
		"/file/:fileName/:anotherParam":         FileNotRandom,
		"/file/:fileName/random/:somethingElse": FileRandom,
	}

	for route, function := range routes {
		testRoute(t, route, function)
	}
}
