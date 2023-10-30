package controller

import (
	"encoding/json"
	"net/http"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			Status: http.StatusMethodNotAllowed,
			Text:   "Method Not Allowed",
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Status: http.StatusOK,
		Text:   "Welcome To My Hello World!",
	})
}
