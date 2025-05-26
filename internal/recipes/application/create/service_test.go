package create

import (
	"context"
	"errors"
	"testing"

	recipesdomain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/recipes/platform/storage/storagemocks"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
	"github.com/rubenbupe/recipe-video-parser/kit/event/eventmocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ExtractionService_CreateExtraction_RepositoryError(t *testing.T) {
	extractionID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	extractionRepositoryMock := new(storagemocks.ExtractionRepository)
	extractionRepositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	extractionRepositoryMock.On("Save", mock.Anything, mock.Anything).Return(errors.New("something unexpected happened"))

	eventBusMock := new(eventmocks.Bus)

	extractionService := NewExtractionService(extractionRepositoryMock, eventBusMock)

	err := extractionService.CreateExtraction(context.Background(), extractionID, userID, data, metadata, createdAt)

	extractionRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_ExtractionService_CreateExtraction_EventsBusError(t *testing.T) {
	extractionID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	extractionRepositoryMock := new(storagemocks.ExtractionRepository)
	extractionRepositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	extractionRepositoryMock.On("Save", mock.Anything, mock.Anything).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(errors.New("something unexpected happened"))

	extractionService := NewExtractionService(extractionRepositoryMock, eventBusMock)

	err := extractionService.CreateExtraction(context.Background(), extractionID, userID, data, metadata, createdAt)

	extractionRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_ExtractionService_CreateExtraction_Succeed(t *testing.T) {
	extractionID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	extractionRepositoryMock := new(storagemocks.ExtractionRepository)
	extractionRepositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
	extractionRepositoryMock.On("Save", mock.Anything, mock.Anything).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.MatchedBy(func(events []event.Event) bool {
		return len(events) > 0
	})).Return(nil)
	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(nil)

	extractionService := NewExtractionService(extractionRepositoryMock, eventBusMock)

	err := extractionService.CreateExtraction(context.Background(), extractionID, userID, data, metadata, createdAt)

	extractionRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.NoError(t, err)
}

func Test_ExtractionService_CreateExtraction_AlreadyExists(t *testing.T) {
	extractionID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	data := "{\"field\":\"value\"}"
	metadata := "{\"meta\":\"value\"}"
	createdAt := "2023-10-01T00:00:00Z"

	extractionRepositoryMock := new(storagemocks.ExtractionRepository)
	extractionRepositoryMock.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

	eventBusMock := new(eventmocks.Bus)

	extractionService := NewExtractionService(extractionRepositoryMock, eventBusMock)

	err := extractionService.CreateExtraction(context.Background(), extractionID, userID, data, metadata, createdAt)

	extractionRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, err, recipesdomain.ErrExtractionAlreadyExists)
}
