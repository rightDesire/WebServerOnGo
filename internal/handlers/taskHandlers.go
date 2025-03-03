package handlers

import (
	"WebServer/internal/errorMessages"
	"WebServer/internal/tasksService"
	"WebServer/internal/web/tasks"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TaskHandler struct {
	Service *tasksService.TaskService
}

// ctor
func NewTaskHandler(taskService *tasksService.TaskService) *TaskHandler {
	return &TaskHandler{
		Service: taskService,
	}
}

func (h *TaskHandler) GetApiTasks(ctx context.Context, request tasks.GetApiTasksRequestObject) (tasks.GetApiTasksResponseObject, error) {
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

func (h *TaskHandler) PostApiTasks(ctx context.Context, request tasks.PostApiTasksRequestObject) (tasks.PostApiTasksResponseObject, error) {
	// Распаковываем тело запроса напрямую, без декодера!
	taskRequest := request.Body
	// Обращаемся к сервису и создаем задачу
	taskToCreate := tasksService.Task{
		UserID: *taskRequest.UserId,
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
		UserId: &createdTask.UserID,
		Task:   &createdTask.Task,
		IsDone: createdTask.IsDone,
	}

	return response, nil
}

func (h *TaskHandler) PatchApiTasksId(ctx context.Context, request tasks.PatchApiTasksIdRequestObject) (tasks.PatchApiTasksIdResponseObject, error) {
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
		if errors.Is(err, errorMessages.ErrNoFieldsToUpdate) {
			errorMsg := "No fields to update"
			return tasks.PatchApiTasksId400JSONResponse{Message: &errorMsg}, nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := "Task not found"
			return tasks.PatchApiTasksId404JSONResponse{Message: &errorMsg}, nil
		}
		if errors.Is(err, echo.ErrBadRequest) {
			errorMsg := "Bad request for update user"
			return tasks.PatchApiTasksId400JSONResponse{Message: &errorMsg}, nil
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

func (h *TaskHandler) DeleteApiTasksId(ctx context.Context, request tasks.DeleteApiTasksIdRequestObject) (tasks.DeleteApiTasksIdResponseObject, error) {
	if err := h.Service.DeleteTaskByID(request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := "Task not found"
			return tasks.DeleteApiTasksId404JSONResponse{Message: &errorMsg}, nil
		}
		return nil, err
	}

	return tasks.DeleteApiTasksId204Response{}, nil
}
