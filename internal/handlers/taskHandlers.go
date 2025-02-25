package handlers

import (
	"WebServer/internal/taskService"
	"WebServer/pkg/utils"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type Handler struct {
	Service *taskService.TaskService
}

// ctor
func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req taskService.Task

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	task, err := h.Service.CreateTask(req)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Database error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&task); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Encoding error")
	}
}

func (h *Handler) GetTasksHandler(w http.ResponseWriter, _ *http.Request) {
	tasks, err := h.Service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Encoding error")
	}
}

func (h *Handler) PatchTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req taskService.Task
	tempId := mux.Vars(r)["id"]
	id, err := utils.StringToInt(tempId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Conversion error")
	}

	// Считываем тело запроса в сырую форму (байты)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Парсим body в map[string]json.RawMessage,
	// чтобы увидеть, какие ключи (поля) пришли
	// json.Unmarshal – распарсит байтовый массив JSON в любой GO-тип
	// json.RawMessage – необработанное значение JSON
	// Использую json.RawMessage, т.к. значения на данном этапе не требуются. Оптимизация декодирования
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal(body, &rawData); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Проверяем, что пришёл ровно один ключ — "is_done"
	if len(rawData) != 1 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Only 'is_done' field is allowed")
		return
	}
	if _, ok := rawData["is_done"]; !ok {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Only 'is_done' field is allowed")
		return
	}

	// Теперь декодируем значение "is_done" в req
	if err := json.Unmarshal(rawData["is_done"], &req.IsDone); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid 'is_done' value")
		return
	}

	task, err := h.Service.UpdateTaskByID(id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.SendErrorResponse(w, http.StatusNotFound, "Task not found")
		default:
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete task")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Encoding error")
	}
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	tempId := mux.Vars(r)["id"]
	id, err := utils.StringToInt(tempId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Conversion error")
	}
	if err := h.Service.DeleteTaskByID(id); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.SendErrorResponse(w, http.StatusNotFound, "Task not found")
		default:
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete task")
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
