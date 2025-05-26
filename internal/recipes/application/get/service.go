package get

import (
	"context"

	extractionsdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
)

type ExtractionService struct {
	extractionRepository extractionsdomain.ExtractionRepository
}

func NewExtractionService(extractionRepository extractionsdomain.ExtractionRepository) ExtractionService {
	return ExtractionService{
		extractionRepository: extractionRepository,
	}
}

func (s ExtractionService) GetExtraction(ctx context.Context, userId string) ([]extractionsdomain.Extraction, error) {
	userID, err := extractionsdomain.NewExtractionUserID(userId)
	if err != nil {
		return nil, err
	}

	extraction, err := s.extractionRepository.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return extraction, nil
}
