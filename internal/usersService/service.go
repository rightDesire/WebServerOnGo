package usersService

import "WebServer/internal/tasksService"

type UserService struct {
	repo *userRepository
}

func NewService(repo *userRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user User) (User, error) { return s.repo.CreateUser(user) }

func (s *UserService) GetAllUsers() ([]User, error) { return s.repo.GetAllUsers() }

func (s *UserService) GetTasksForUser(userID uint, tasks []tasksService.Task) ([]tasksService.Task, error) {
	var response []tasksService.Task
	for _, task := range tasks {
		if task.ID == userID {
			response = append(response, task)
		}
	}
	return response, nil
}

func (s *UserService) UpdateUserByID(id uint, user User) (User, error) {
	return s.repo.UpdateUserByID(id, user)
}

func (s *UserService) DeleteUserByID(id uint) error { return s.repo.DeleteUserByID(id) }
