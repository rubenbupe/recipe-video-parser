package get

import (
	"context"
	"errors"
	"testing"

	recipesdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/storage/storagemocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ExtractionService_GetExtraction_RepositoryError(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"

	extractionRepositoryMock := new(storagemocks.ExtractionRepository)
	extractionRepositoryMock.On("GetByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("something unexpected happened"))

	extractionService := NewExtractionService(extractionRepositoryMock)

	_, err := extractionService.GetExtraction(context.Background(), userID)

	extractionRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_ExtractionService_GetExtraction_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	id := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	extraction, err := recipesdomain.NewExtraction(id, userID, data, metadata, createdAt)
	assert.NoError(t, err)

	extractionRepositoryMock := new(storagemocks.ExtractionRepository)
	extractionRepositoryMock.On("GetByUserID", mock.Anything, mock.Anything).Return([]recipesdomain.Extraction{extraction}, nil)

	extractionService := NewExtractionService(extractionRepositoryMock)

	foundExtractions, err := extractionService.GetExtraction(context.Background(), userID)

	assert.Equal(t, 1, len(foundExtractions))
	assert.Equal(t, extraction.Id.String(), foundExtractions[0].Id.String())
	extractionRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
