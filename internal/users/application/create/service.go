package create

import (
	"context"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

type UserService struct {
	userRepository usersdomain.UserRepository
	eventBus       event.Bus
}

func NewUserService(userRepository usersdomain.UserRepository, eventBus event.Bus) UserService {
	return UserService{
		userRepository: userRepository,
		eventBus:       eventBus,
	}
}

func (s UserService) CreateUser(ctx context.Context, id, name, apikey, createdAt string) error {
	userID, err := usersdomain.NewUserID(id)
	if err != nil {
		return err
	}

	userName, err := usersdomain.NewUserName(name)
	if err != nil {
		return err
	}

	userExists, err := s.userRepository.Exists(ctx, userID)
	if err != nil {
		return err
	}

	userExistsByName, err := s.userRepository.ExistsByName(ctx, userName)

	if err != nil {
		return err
	}

	if userExists || userExistsByName {
		return usersdomain.ErrUserAlreadyExists
	}

	user, err := usersdomain.NewUser(id, name, apikey, createdAt)
	if err != nil {
		return err
	}

	if err := s.userRepository.Save(ctx, user); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, user.PullEvents())
}
