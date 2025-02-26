package handlers

import (
	"WebServer/internal/tasksService"
	"WebServer/internal/web/tasks"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Handler struct {
	Service *tasksService.TaskService
}

// ctor
func NewHandler(service *tasksService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetApiTasks(ctx context.Context, request tasks.GetApiTasksRequestObject) (tasks.GetApiTasksResponseObject, error) {
	// Получение всех задач из сервиса
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	// Создаем переменную, которую мы потом передадим в качестве ответа
	response := tasks.GetApiTasks200JSONResponse{} // Заполняем слайс response всеми задачами из БД
	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     &tsk.ID,
			Task:   &tsk.Task,
			IsDone: tsk.IsDone,
		}
		response = append(response, task)
	}

	return response, nil
}

func (h *Handler) PostApiTasks(ctx context.Context, request tasks.PostApiTasksRequestObject) (tasks.PostApiTasksResponseObject, error) {
	// Распаковываем тело запроса напрямую, без декодера!
	taskRequest := request.Body
	// Обращаемся к сервису и создаем задачу
	taskToCreate := tasksService.Task{
		Task:   *taskRequest.Task,
		IsDone: taskRequest.IsDone,
	}
	createdTask, err := h.Service.CreateTask(taskToCreate)

	if err != nil {
		return nil, err
	}
	// создаем структуру ответа
	response := tasks.PostApiTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Task,
		IsDone: createdTask.IsDone,
	}

	return response, nil
}

func (h *Handler) PatchApiTasksId(ctx context.Context, request tasks.PatchApiTasksIdRequestObject) (tasks.PatchApiTasksIdResponseObject, error) {
	taskRequest := request.Body
	taskToUpdate := tasksService.Task{}

	if taskRequest.Task != nil {
		taskToUpdate.Task = *taskRequest.Task
	}
	if taskRequest.IsDone != nil {
		taskToUpdate.IsDone = taskRequest.IsDone
	}

	task, err := h.Service.UpdateTaskByID(request.Id, taskToUpdate)
	if err != nil {
		if errors.Is(err, tasksService.ErrNoFieldsToUpdate) {
			errorMsg := "No fields to update"
			return tasks.PatchApiTasksId400JSONResponse{Message: &errorMsg}, nil
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.PatchApiTasksId404Response{}, nil
		}
		return nil, err
	}

	response := tasks.PatchApiTasksId200JSONResponse{
		Id:     &task.ID,
		Task:   &task.Task,
		IsDone: task.IsDone,
	}

	return response, nil
}

func (h *Handler) DeleteApiTasksId(ctx context.Context, request tasks.DeleteApiTasksIdRequestObject) (tasks.DeleteApiTasksIdResponseObject, error) {
	if err := h.Service.DeleteTaskByID(request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks.DeleteApiTasksId404Response{}, nil
		}
		return nil, err
	}

	return tasks.DeleteApiTasksId204Response{}, nil
}
