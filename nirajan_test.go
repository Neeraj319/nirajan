package nirajan

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimpleRouter(t *testing.T) {
	router := CreateRouter()

	router.AddRoute("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, GET)

	type Params struct {
		Name string
		Age  int
	}
	router.AddRoute("/:Name/:Age", func(w http.ResponseWriter, r *http.Request, params Params) {
		w.WriteHeader(http.StatusCreated)
	}, GET)

	router.AddRoute("/users/:Id", func(w http.ResponseWriter, r *http.Request, params struct{ Id string }) {
		if params.Id != "1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}, GET)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/users/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	req, err = http.NewRequest("GET", "/users/geda", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
