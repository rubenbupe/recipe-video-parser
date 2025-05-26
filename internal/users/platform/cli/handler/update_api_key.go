package handlers

import (
	"context"
	"fmt"

	"github.com/rubenbupe/recipe-video-parser/internal/users/application/updateapikey"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

type UpdateApiKeyInput struct {
	Name   string
	ApiKey string
}

type UpdateApiKeyHandler func(context.Context, UpdateApiKeyInput) error

func CreateUpdateApiKeyHandler(commandBus command.Bus) UpdateApiKeyHandler {
	return func(ctx context.Context, input UpdateApiKeyInput) error {
		if input.Name == "" || input.ApiKey == "" {
			return fmt.Errorf("todos los campos son obligatorios")
		}

		err := commandBus.Dispatch(ctx, updateapikey.NewApiKeyUserCommand(
			input.Name,
			input.ApiKey,
		))

		if err != nil {
			switch {
			case err == usersdomain.ErrEmptyUserName:
				return fmt.Errorf("error de dominio: %w", err)
			default:
				return fmt.Errorf("error interno: %w", err)
			}
		}

		return nil
	}
}
