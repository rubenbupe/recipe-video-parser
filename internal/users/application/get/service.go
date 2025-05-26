package get

import (
	"context"
	"fmt"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
)

type UserService struct {
	userRepository usersdomain.UserRepository
}

func NewUserService(userRepository usersdomain.UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (s UserService) GetUser(ctx context.Context, name string) (*usersdomain.User, error) {
	userName, err := usersdomain.NewUserName(name)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetByName(ctx, userName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}
	return user, nil
}
