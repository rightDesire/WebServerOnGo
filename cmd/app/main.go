package main

import (
	"WebServer/internal/database"
	"WebServer/internal/handlers"
	"WebServer/internal/taskService"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	database.InitDB()

	repo := taskService.NewTaskRepository(database.DB)
	service := taskService.NewService(repo)

	handler := handlers.NewHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", handler.PostTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/tasks", handler.GetTasksHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/tasks/{id}", handler.PatchTaskHandler).Methods(http.MethodPatch)
	router.HandleFunc("/api/tasks/{id}", handler.DeleteTaskHandler).Methods(http.MethodDelete)
	http.ListenAndServe(":8080", router)
}
