package handlers

import (
	"context"
	"fmt"

	"github.com/rubenbupe/recipe-video-parser/internal/recipes/application/get"
	recipesdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

type GetExtractionInput struct {
	UserID string
}

type GetExtractionOutput struct {
	ID        string
	UserID    string
	Data      string
	Metadata  string
	CreatedAt string
}

type GetExtractionHandler func(context.Context, GetExtractionInput) ([]GetExtractionOutput, error)

func CreateGetExtractionsHandler(queryBus query.Bus) GetExtractionHandler {
	return func(ctx context.Context, input GetExtractionInput) ([]GetExtractionOutput, error) {
		if input.UserID == "" {
			return nil, fmt.Errorf("el campo ID es obligatorio")
		}

		result, err := queryBus.Ask(ctx, get.NewExtractionQuery(input.UserID))
		if err != nil {
			return nil, fmt.Errorf("error al buscar extracciones: %w", err)
		}

		extractions, ok := result.([]recipesdomain.Extraction)
		if !ok {
			return nil, fmt.Errorf("error al convertir el resultado a tipo []Extraction")
		}

		if len(extractions) == 0 {
			return nil, fmt.Errorf("no se encontraron extracciones para el usuario con ID: %s", input.UserID)
		}

		outputs := make([]GetExtractionOutput, 0, len(extractions))
		for _, e := range extractions {
			outputs = append(outputs, GetExtractionOutput{
				ID:        e.Id.String(),
				UserID:    e.UserId.String(),
				Data:      e.Data,
				Metadata:  e.Metadata,
				CreatedAt: e.CreatedAt.String(),
			})
		}
		return outputs, nil
	}
}
