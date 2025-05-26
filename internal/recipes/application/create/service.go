package create

import (
	"context"

	extractionsdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

type ExtractionService struct {
	extractionRepository extractionsdomain.ExtractionRepository
	eventBus       event.Bus
}

func NewExtractionService(extractionRepository extractionsdomain.ExtractionRepository, eventBus event.Bus) ExtractionService {
	return ExtractionService{
		extractionRepository: extractionRepository,
		eventBus:       eventBus,
	}
}

func (s ExtractionService) CreateExtraction(ctx context.Context, id, userId, data, metadata, createdAt string) error {
	extractionId, err := extractionsdomain.NewExtractionID(id)
	if err != nil {
		return err
	}

	extractionExists, err := s.extractionRepository.Exists(ctx, extractionId)

	if err != nil {
		return err
	}

	if extractionExists {
		return extractionsdomain.ErrExtractionAlreadyExists
	}

	extraction, err := extractionsdomain.NewExtraction(id, userId, data, metadata, createdAt)
	if err != nil {
		return err
	}

	if err := s.extractionRepository.Save(ctx, extraction); err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, extraction.PullEvents())
}
