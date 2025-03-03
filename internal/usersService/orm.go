package usersService

import (
	"WebServer/internal/tasksService"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string              `json:"email"`
	Password string              `json:"password"`
	Tasks    []tasksService.Task `gorm:"foreignKey:UserID"`
}
