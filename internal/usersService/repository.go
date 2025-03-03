package usersService

import (
	"WebServer/internal/errorMessages"
	"WebServer/internal/tasksService"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
	GetTasksByUserID(userID uint) (User, error)
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

func (repo *userRepository) GetTasksByUserID(userID uint) (User, error) {
	// Определяем "промежуточную" структуру, в которую будем сканировать результат JOIN.
	type userTaskRow struct {
		UserID   uint   `gorm:"column:user_id"`
		Email    string `gorm:"column:email"`
		Password string `gorm:"column:password"`

		TaskID *uint   `gorm:"column:task_id"`
		Task   *string `gorm:"column:task"`
		IsDone *bool   `gorm:"column:is_done"`
	}

	var rows []userTaskRow
	res := repo.db.Table("users").
		Select(`
            users.id        AS user_id,
            users.email     AS email,
            users.password  AS password,
            tasks.id        AS task_id,
            tasks.task      AS task,
            tasks.is_done   AS is_done
        `).
		Joins("LEFT JOIN tasks ON tasks.user_id = users.id").
		Where("users.id = ?", userID).
		Scan(&rows)
	if res.Error != nil {
		return User{}, res.Error
	}

	// Если ничего не вернулось, значит пользователя с таким id нет
	if len(rows) == 0 {
		return User{}, gorm.ErrRecordNotFound
	}

	// Создаем объект User из первой строки, т.к. данные о пользователе у всех строк будут один и те же
	user := User{
		Model: gorm.Model{
			ID: rows[0].UserID,
		},
		Email:    rows[0].Email,
		Password: rows[0].Password,
		Tasks:    []tasksService.Task{},
	}

	// Перебираем все строки и формируем массив Tasks
	for _, row := range rows {
		// Если в строке нет TaskID (NULL в базе),
		// значит у этого пользователя пока нет связанной задачи
		if row.TaskID != nil {
			user.Tasks = append(user.Tasks, tasksService.Task{
				Model: gorm.Model{
					ID: *row.TaskID,
				},
				Task:   *row.Task,
				IsDone: row.IsDone,
				UserID: row.UserID,
			})
		}
	}

	return user, nil
}
