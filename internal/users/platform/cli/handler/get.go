package handlers

import (
	"context"
	"fmt"

	"github.com/rubenbupe/recipe-video-parser/internal/users/application/get"
	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

type GetUserInput struct {
	Name string
}

type GetUserOutput struct {
	ID        string
	Name      string
	ApiKey    string
	CreatedAt string
}

type GetUserHandler func(context.Context, GetUserInput) (*GetUserOutput, error)

func CreateGetUserHandler(queryBus query.Bus) GetUserHandler {
	return func(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
		if input.Name == "" {
			return nil, fmt.Errorf("el campo ID es obligatorio")
		}

		result, err := queryBus.Ask(ctx, get.NewUserQuery(input.Name))
		if err != nil {
			return nil, fmt.Errorf("error al buscar usuario: %w", err)
		}

		user, ok := result.(*usersdomain.User)
		if !ok {
			return nil, fmt.Errorf("respuesta inesperada del query")
		}

		return &GetUserOutput{
			ID:        user.Id.String(),
			Name:      user.Name.String(),
			ApiKey:    user.ApiKey.String(),
			CreatedAt: user.CreatedAt.String(),
		}, nil
	}
}
