package updateapikey

import (
	"context"
	"errors"
	"testing"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/storagemocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserService_UpdateApiKey_RepositoryError(t *testing.T) {
	userName := "Test User"
	newApiKey := "new-api-key"

	repo := new(storagemocks.UserRepository)
	repo.On("GetByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(nil, errors.New("not found"))

	service := NewUserApiKeyService(repo, nil)
	err := service.UpdateApiKey(context.Background(), userName, newApiKey)

	repo.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_UpdateApiKey_SaveError(t *testing.T) {
	userName := "Test User"
	newApiKey := "new-api-key"
	user, _ := usersdomain.NewUser("37a0f027-15e6-47cc-a5d2-64183281087e", userName, "old-api-key", "2023-10-01T00:00:00Z")

	repo := new(storagemocks.UserRepository)
	repo.On("GetByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(&user, nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("domain.User")).Return(errors.New("save error"))

	service := NewUserApiKeyService(repo, nil)
	err := service.UpdateApiKey(context.Background(), userName, newApiKey)

	repo.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_UpdateApiKey_Success(t *testing.T) {
	userName := "Test User"
	newApiKey := "new-api-key"
	user, _ := usersdomain.NewUser("37a0f027-15e6-47cc-a5d2-64183281087e", userName, "old-api-key", "2023-10-01T00:00:00Z")

	repo := new(storagemocks.UserRepository)
	repo.On("GetByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(&user, nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)

	service := NewUserApiKeyService(repo, nil)
	err := service.UpdateApiKey(context.Background(), userName, newApiKey)

	repo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, newApiKey, user.ApiKey.String())
}
