package handlers

import (
	"WebServer/internal/errorMessages"
	"WebServer/internal/usersService"
	"WebServer/internal/web/users"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	Service *usersService.UserService
}

func NewUserHandler(service *usersService.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) GetApiUsersUserIdTasks(ctx context.Context, request users.GetApiUsersUserIdTasksRequestObject) (users.GetApiUsersUserIdTasksResponseObject, error) {
	tasks, err := h.Service.GetTasksForUser(request.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := "User not found"
			return users.GetApiUsersUserIdTasks404JSONResponse{Message: &errorMsg}, nil
		}
		return nil, err
	}

	var response users.GetApiUsersUserIdTasks200JSONResponse
	for _, tsk := range tasks {
		task := users.Task{
			Id:     &tsk.ID,
			Task:   &tsk.Task,
			IsDone: tsk.IsDone,
			UserId: &tsk.UserID,
		}
		response = append(response, task)
	}
	return response, err
}

func (h UserHandler) GetApiUsers(ctx context.Context, _ users.GetApiUsersRequestObject) (users.GetApiUsersResponseObject, error) {
	data, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := users.GetApiUsers200JSONResponse{}
	for _, usr := range data {
		user := users.User{
			Id:       &usr.ID,
			Email:    &usr.Email,
			Password: &usr.Password,
		}
		response = append(response, user)
	}

	return response, nil
}

func (h UserHandler) PostApiUsers(ctx context.Context, request users.PostApiUsersRequestObject) (users.PostApiUsersResponseObject, error) {
	userRequest := request.Body
	userToCreate := usersService.User{
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
	}
	data, err := h.Service.CreateUser(userToCreate)
	if err != nil {
		return nil, err
	}

	response := &users.PostApiUsers201JSONResponse{
		Id:       &data.ID,
		Email:    &data.Email,
		Password: &data.Password,
	}
	return response, nil
}

func (h UserHandler) DeleteApiUsersId(ctx context.Context, request users.DeleteApiUsersIdRequestObject) (users.DeleteApiUsersIdResponseObject, error) {
	if err := h.Service.DeleteUserByID(request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := "User not found"
			return users.DeleteApiUsersId404JSONResponse{Message: &errorMsg}, nil
		}
		return nil, err
	}
	return users.DeleteApiUsersId204Response{}, nil
}

func (h UserHandler) PatchApiUsersId(ctx context.Context, request users.PatchApiUsersIdRequestObject) (users.PatchApiUsersIdResponseObject, error) {
	userRequest := request.Body
	userToUpdate := usersService.User{}

	if userRequest.Email != nil {
		userToUpdate.Email = *userRequest.Email
	}
	if userRequest.Password != nil {
		userToUpdate.Password = *userRequest.Password
	}

	data, err := h.Service.UpdateUserByID(request.Id, userToUpdate)
	if err != nil {
		if errors.Is(err, errorMessages.ErrNoFieldsToUpdate) {
			errorMsg := "No fields to update"
			return users.PatchApiUsersId400JSONResponse{Message: &errorMsg}, nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := "User not found"
			return users.PatchApiUsersId404JSONResponse{Message: &errorMsg}, nil
		}
		if errors.Is(err, echo.ErrBadRequest) {
			errorMsg := "Bad request for update user"
			return users.PatchApiUsersId400JSONResponse{Message: &errorMsg}, nil
		}
		return nil, err
	}

	response := &users.PatchApiUsersId200JSONResponse{
		Id:       &data.ID,
		Email:    &data.Email,
		Password: &data.Password,
	}
	return response, nil
}
