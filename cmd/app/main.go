package main

import (
	"WebServer/internal/database"
	"WebServer/internal/handlers"
	"WebServer/internal/tasksService"
	"WebServer/internal/usersService"
	"WebServer/internal/web/tasks"
	"WebServer/internal/web/users"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	database.InitDB()

	tasksRepo := tasksService.NewRepository(database.DB)
	tasksServ := tasksService.NewService(tasksRepo)
	tasksHandler := handlers.NewTaskHandler(tasksServ)

	usersRepo := usersService.NewRepository(database.DB)
	usersServ := usersService.NewService(usersRepo)
	usersHandler := handlers.NewUserHandler(usersServ)

	// Инициализируем echo
	e := echo.New()

	// используем Logger и Recover
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрируем хэндлеры в echo
	tasksStrictHandler := tasks.NewStrictHandler(tasksHandler, nil) // тут будет ошибка
	tasks.RegisterHandlers(e, tasksStrictHandler)
	usersStrictHandler := users.NewStrictHandler(usersHandler, nil)
	users.RegisterHandlers(e, usersStrictHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
