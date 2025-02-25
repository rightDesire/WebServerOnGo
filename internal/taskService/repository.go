package taskService

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTasks() ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

// ctor
func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	result := r.db.Model(&Task{}).Where("ID = ?", id).Update("is_done", task.IsDone)
	if result.Error != nil {
		return Task{}, result.Error
	}
	if err := r.db.Where("ID = ?", id).First(&task).Error; err != nil {
		return Task{}, err
	}
	return task, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	result := r.db.Model(&Task{}).Where("ID = ?", id).Delete(&Task{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
