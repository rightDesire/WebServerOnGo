package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskStruct Task

	if err := json.NewDecoder(r.Body).Decode(&taskStruct); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result := DB.Create(&taskStruct)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&taskStruct); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

func GetTask(w http.ResponseWriter, _ *http.Request) {
	var tasks []Task
	DB.Find(&tasks)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

func main() {
	InitDB()

	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks", GetTask).Methods("GET")
	http.ListenAndServe(":8080", router)
}
