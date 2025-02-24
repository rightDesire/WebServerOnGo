package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var req Task

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result := DB.Create(&req)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&req); err != nil {
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

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Анонимная структура
	var req struct {
		IsDone bool `json:"is_done"`
	}

	var task Task
	id := mux.Vars(r)["id"]

	// Считываем тело запроса в сырую форму (байты)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Парсим body в map[string]json.RawMessage,
	// чтобы увидеть, какие ключи (поля) пришли
	// json.Unmarshal – распарсит байтовый массив JSON в любой GO-тип
	// json.RawMessage – необработанное значение JSON
	// Использую json.RawMessage, т.к. значения на данном этапе не требуются. Оптимизация декодирования
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal(body, &rawData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что пришёл ровно один ключ — "is_done"
	if len(rawData) != 1 {
		http.Error(w, "Only 'is_done' field is allowed", http.StatusBadRequest)
		return
	}
	if _, ok := rawData["is_done"]; !ok {
		http.Error(w, "Only 'is_done' field is allowed", http.StatusBadRequest)
		return
	}

	// Теперь декодируем значение "is_done" в req
	if err := json.Unmarshal(rawData["is_done"], &req.IsDone); err != nil {
		http.Error(w, "Invalid 'is_done' value", http.StatusBadRequest)
		return
	}

	result := DB.Model(&Task{}).Where("ID = ?", id).Update("is_done", req.IsDone)
	if result.Error != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	if err := DB.Where("ID = ?", id).First(&task).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result := DB.Model(&Task{}).Where("ID = ?", id).Delete(&Task{})
	if result.Error != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/api/tasks", GetTask).Methods(http.MethodGet)
	router.HandleFunc("/api/tasks/{id}", UpdateTask).Methods(http.MethodPatch)
	router.HandleFunc("/api/tasks/{id}", DeleteTask).Methods(http.MethodDelete)
	http.ListenAndServe(":8080", router)
}
