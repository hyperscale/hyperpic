package main

import "net/http"

// HTTPFailure make http error
func HTTPFailure(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
}
