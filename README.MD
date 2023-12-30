# Nirajan

Another router for the net/http module 


## Get Started

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/Neeraj319/nirajan"
)

type Params struct {
	Name string // won't work if atribute name does not start with a Capital letter
	Age  int
}

func Home(w http.ResponseWriter, r *http.Request, params Params) {
	response_string := fmt.Sprintf("name: %s age: %d \n", params.Name, params.Age)
	fmt.Fprintf(w, response_string)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println(nirajan.QueryParams(r))
	fmt.Fprintf(w, "index response")
}

func main() {
	r := nirajan.CreateRouter()
	r.AddRoute("/:Name/:Age", Home, nirajan.GET)
	r.AddRoute("/", Index, nirajan.GET)
	http.ListenAndServe("0.0.0.0:8080", r)

}
```

- in the memory of nirajan dhakal (he is not dead)
