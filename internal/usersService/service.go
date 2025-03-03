package usersService

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

func (s *UserService) GetTasksForUser(userID uint) (User, error) {
	return s.repo.GetTasksByUserId(userID)
}

func (s *UserService) UpdateUserByID(id uint, user User) (User, error) {
	return s.repo.UpdateUserByID(id, user)
}

func (s *UserService) DeleteUserByID(id uint) error { return s.repo.DeleteUserByID(id) }
