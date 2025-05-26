package getbyapikey

import (
	"context"
	"errors"
	"testing"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/storagemocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserByApiKeyService_GetUserByApiKey_RepositoryError(t *testing.T) {
	apiKey := "test-api-key"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("GetByApiKey", mock.Anything, mock.Anything).Return(((*usersdomain.User)(nil)), errors.New("something unexpected happened"))

	userService := NewUserByApiKeyService(userRepositoryMock)

	_, err := userService.GetUser(context.Background(), apiKey)

	userRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserByApiKeyService_GetUserByApiKey_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
	userApiKey := "test-api-key"
	userCreatedAt := "2023-10-01T00:00:00Z"

	user, err := usersdomain.NewUser(userID, userName, userApiKey, userCreatedAt)
	assert.NoError(t, err)

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("GetByApiKey", mock.Anything, mock.Anything).Return(&user, nil)

	userService := NewUserByApiKeyService(userRepositoryMock)

	foundUser, err := userService.GetUser(context.Background(), userApiKey)

	assert.Equal(t, user.Name.String(), foundUser.Name.String())
	userRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
