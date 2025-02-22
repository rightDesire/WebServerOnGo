package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskStruct Task

	err := json.NewDecoder(r.Body).Decode(&taskStruct)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	DB.Create(&taskStruct)

	err = json.NewEncoder(w).Encode(&taskStruct)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
}

func GetTask(w http.ResponseWriter, _ *http.Request) {
	var tasks []Task
	DB.Find(&tasks)

	err := json.NewEncoder(w).Encode(&tasks)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
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
