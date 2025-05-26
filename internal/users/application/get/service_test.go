package get

import (
	"context"
	"errors"
	"testing"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/storagemocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserService_GetUser_RepositoryError(t *testing.T) {
	userName := "Test User"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("GetByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(((*usersdomain.User)(nil)), errors.New("something unexpected happened"))

	userService := NewUserService(userRepositoryMock)

	_, err := userService.GetUser(context.Background(), userName)

	userRepositoryMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_GetUser_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
  userApiKey := "test-api-key"
  userCreatedAt := "2023-10-01T00:00:00Z"

	user, err := usersdomain.NewUser(userID, userName, userApiKey, userCreatedAt)
	assert.NoError(t, err)

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("GetByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(&user, nil)

	userService := NewUserService(userRepositoryMock)

	foundUser, err := userService.GetUser(context.Background(), userName)

	assert.Equal(t, user.Name.String(), foundUser.Name.String())
	userRepositoryMock.AssertExpectations(t)
	assert.NoError(t, err)
}
