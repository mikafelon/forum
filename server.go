package main

import (
	"fmt"
	"net/http"
)

// launch a server on port 8080 and display auth.html
func server() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	fmt.Println("Server is running on port 8080.")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
