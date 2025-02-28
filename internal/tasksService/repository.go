package tasksService

import (
	"errors"
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
func NewRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

var ErrNoFieldsToUpdate = errors.New("no fields to update")

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
	var updateTask Task
	updateData := make(map[string]interface{})

	if task.Task != "" {
		updateData["task"] = task.Task
	}
	if task.IsDone != nil {
		updateData["is_done"] = task.IsDone
	}

	if len(updateData) == 0 {
		return Task{}, ErrNoFieldsToUpdate
	}

	result := r.db.Model(&Task{}).Where("ID = ?", id).Updates(updateData)
	if result.Error != nil {
		return Task{}, result.Error
	}

	if result.RowsAffected == 0 {
		return Task{}, gorm.ErrRecordNotFound
	}

	result = r.db.Where("ID = ?", id).First(&updateTask)
	if result.Error != nil {
		return Task{}, result.Error
	}

	return updateTask, nil
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	var task Task
	result := r.db.Model(&Task{}).Where("ID = ?", id).Delete(&task)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
