package getbyapikey

import (
	"context"
	"fmt"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
)

type UserByApiKeyService struct {
	userRepository usersdomain.UserRepository
}

func NewUserByApiKeyService(userRepository usersdomain.UserRepository) UserByApiKeyService {
	return UserByApiKeyService{
		userRepository: userRepository,
	}
}

func (s UserByApiKeyService) GetUser(ctx context.Context, apiKey string) (*usersdomain.User, error) {
	userApiKey, err := usersdomain.NewUserApiKey(apiKey)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepository.GetByApiKey(ctx, userApiKey)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}
	return user, nil
}
