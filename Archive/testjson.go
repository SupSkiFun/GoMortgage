package main

import (
	"encoding/json"
	"net/http"
)

// JSONResponse works with test() to provide a simple JSON response.
type JSONResponse struct {
	Status string `json:"status"`
}

// test provides a basic JSON response.
func test(w http.ResponseWriter, r *http.Request) {
	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)
}
