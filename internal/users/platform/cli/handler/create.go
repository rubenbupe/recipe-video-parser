package handlers

import (
	"context"
	"fmt"

	"github.com/rubenbupe/recipe-video-parser/internal/users/application/create"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

type CreateUserInput struct {
	ID        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ApiKey    string `json:"apikey" binding:"required"`
	CreatedAt string `json:"createdAt" binding:"required"`
}

type CreateUserHandler func(context.Context, CreateUserInput) error

func CreateHandler(commandBus command.Bus) CreateUserHandler {
	return func(ctx context.Context, input CreateUserInput) error {
		if input.ID == "" || input.Name == "" || input.ApiKey == "" || input.CreatedAt == "" {
			return fmt.Errorf("todos los campos son obligatorios")
		}

		err := commandBus.Dispatch(ctx, create.NewUserCommand(
			input.ID,
			input.Name,
			input.ApiKey,
			input.CreatedAt,
		))

		if err != nil {
			switch {
			case err == usersdomain.ErrInvalidUserID,
				err == usersdomain.ErrEmptyUserName,
				err == usersdomain.ErrUserAlreadyExists:
				return fmt.Errorf("error de dominio: %w", err)
			default:
				return fmt.Errorf("error interno: %w", err)
			}
		}

		return nil
	}

}
