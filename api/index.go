package handler

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api":
		w.Write([]byte("Home Route"))
	case "/api/users":
		w.Write([]byte("User List"))
	default:
		http.NotFound(w, r)
	}
}
