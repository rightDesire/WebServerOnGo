package usersService

import (
	"WebServer/internal/errorMessages"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
	GetTasksByUserId(userID uint) (User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) CreateUser(user User) (User, error) {
	if res := repo.db.Create(&user); res.Error != nil {
		return User{}, res.Error
	}
	return user, nil
}

func (repo *userRepository) GetAllUsers() ([]User, error) {
	var data []User
	if res := repo.db.Model(&User{}).Find(&data); res.Error != nil {
		return []User{}, res.Error
	}

	return data, nil
}

func (repo *userRepository) UpdateUserByID(id uint, user User) (User, error) {
	updateData := make(map[string]interface{})
	if user.Email != "" {
		updateData["email"] = user.Email
	}
	if user.Password != "" {
		updateData["password"] = user.Password
	}
	if len(updateData) == 0 {
		return User{}, errorMessages.ErrNoFieldsToUpdate
	}

	if res := repo.db.Model(&User{}).Where("id = ?", id).Updates(&updateData); res.Error != nil {
		if res.RowsAffected == 0 {
			return User{}, gorm.ErrRecordNotFound
		}
		return User{}, res.Error
	}

	var updateUser User
	res := repo.db.Where("ID = ?", id).First(&updateUser)
	if res.Error != nil {
		return User{}, res.Error
	}

	return updateUser, nil
}

func (repo *userRepository) DeleteUserByID(id uint) error {
	if res := repo.db.Where("id = ?", id).Delete(&User{}); res.Error != nil {
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return res.Error
	}
	return nil
}

func (repo *userRepository) GetTasksByUserId(userID uint) (User, error) {
	var user User
	if res := repo.db.Preload("Tasks").First(&user, userID); res.Error != nil {
		if res.RowsAffected == 0 {
			return User{}, gorm.ErrRecordNotFound
		}
		return User{}, res.Error
	}
	return user, nil
}
