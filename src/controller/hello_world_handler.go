package controller

import (
	"encoding/json"
	"net/http"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
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
