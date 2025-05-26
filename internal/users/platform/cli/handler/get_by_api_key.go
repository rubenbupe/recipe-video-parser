package handlers

import (
	"context"
	"fmt"

	"github.com/rubenbupe/recipe-video-parser/internal/users/application/getbyapikey"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

type GetUserByApiKeyInput struct {
	ApiKey string
}

type GetUserByApiKeyOutput struct {
	ID        string
	Name      string
	ApiKey    string
	CreatedAt string
}

type GetUserByApiKeyHandler func(context.Context, GetUserByApiKeyInput) (*GetUserByApiKeyOutput, error)

func CreateGetUserByApiKeyHandler(queryBus query.Bus) GetUserByApiKeyHandler {
	return func(ctx context.Context, input GetUserByApiKeyInput) (*GetUserByApiKeyOutput, error) {
		if input.ApiKey == "" {
			return nil, fmt.Errorf("el campo ApiKey es obligatorio")
		}

		result, err := queryBus.Ask(ctx, getbyapikey.NewUserQueryByApiKey(input.ApiKey))
		if err != nil {
			return nil, fmt.Errorf("error al buscar usuario por api key: %w", err)
		}

		user, ok := result.(*usersdomain.User)
		if !ok {
			return nil, fmt.Errorf("respuesta inesperada del query")
		}

		return &GetUserByApiKeyOutput{
			ID:        user.Id.String(),
			Name:      user.Name.String(),
			ApiKey:    user.ApiKey.String(),
			CreatedAt: user.CreatedAt.String(),
		}, nil
	}
}
