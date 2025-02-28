package usersService

import "gorm.io/gorm"

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
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

}

func (repo *userRepository) GetAllUsers() ([]User, error) {

}

func (repo *userRepository) UpdateUserByID(id uint, user User) (User, error) {

}

func (repo *userRepository) DeleteUserByID(id uint) error {

}
