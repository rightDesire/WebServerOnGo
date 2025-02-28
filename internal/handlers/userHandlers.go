package handlers

import (
	"WebServer/internal/usersService"
	"WebServer/internal/web/users"
	"context"
)

type UserHandler struct {
	Service *usersService.UserService
}

func NewUserHandler(service *usersService.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (u UserHandler) GetApiUsers(ctx context.Context, request users.GetApiUsersRequestObject) (users.GetApiUsersResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) PostApiUsers(ctx context.Context, request users.PostApiUsersRequestObject) (users.PostApiUsersResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) DeleteApiUsersId(ctx context.Context, request users.DeleteApiUsersIdRequestObject) (users.DeleteApiUsersIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) PatchApiUsersId(ctx context.Context, request users.PatchApiUsersIdRequestObject) (users.PatchApiUsersIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
