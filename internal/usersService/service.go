package usersService

type UserService struct {
	repo *userRepository
}

func NewService(repo *userRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user User) (User, error) {

}

func (s *UserService) GetAllUsers() ([]User, error) {

}

func (s *UserService) UpdateUserByID(id uint, user User) (User, error) {

}

func (s *UserService) DeleteUserByID(id uint) error {

}
