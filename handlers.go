package main

import (
	"fmt"
	"net/http"
)

func v1(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	foo := req.URL.Query()

	for k, v := range foo {
		fmt.Printf("k/v: %v/%v\n", k, v)
	}
}

// BUG(high) gross
func errorHandler(w http.ResponseWriter, req *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404")
	}

	fmt.Println(req)
}

// BUG(high) gross
func catchall(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "welcome home")
}
