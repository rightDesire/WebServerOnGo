package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type requestBody struct {
	Message string `json:"message"`
}

var task string

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var req requestBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task = req.Message
}

func GetHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello,", task)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/task", GetHandler).Methods("GET")
	router.HandleFunc("/api/task", PostHandler).Methods("POST")
	http.ListenAndServe(":8080", router)
}
