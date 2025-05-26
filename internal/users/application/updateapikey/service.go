package updateapikey

import (
	"context"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

type UserApiKeyService struct {
	userRepository usersdomain.UserRepository
	eventBus       event.Bus
}

func NewUserApiKeyService(userRepository usersdomain.UserRepository, eventBus event.Bus) UserApiKeyService {
	return UserApiKeyService{
		userRepository: userRepository,
		eventBus:       eventBus,
	}
}

func (s UserApiKeyService) UpdateApiKey(ctx context.Context, name, apikey string) error {
	userName, err := usersdomain.NewUserName(name)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetByName(ctx, userName)
	if err != nil {
		return err
	}
	if user == nil {
		return usersdomain.ErrEmptyUserName // o un error más específico de usuario no encontrado
	}

	err = user.SetApiKey(apikey)
	if err != nil {
		return err
	}

	if err := s.userRepository.Save(ctx, *user); err != nil {
		return err
	}

	if s.eventBus != nil {
		return s.eventBus.Publish(ctx, user.PullEvents())
	}
	return nil
}
