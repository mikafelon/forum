package main

import (
	"net/http"
)

// launch a server on port 8080 and display auth.html
func server() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe("localhost:8080", nil)
}
