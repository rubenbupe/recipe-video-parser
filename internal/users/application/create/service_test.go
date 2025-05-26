package create

import (
	"context"
	"errors"
	"testing"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/internal/users/platform/storage/storagemocks"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
	"github.com/rubenbupe/recipe-video-parser/kit/event/eventmocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserService_CreateUser_RepositoryError(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
  userApiKey := "test-api-key"
  userCreatedAt := "2023-10-01T00:00:00Z"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Exists", mock.Anything, mock.AnythingOfType("domain.UserID")).Return(false, nil)
  userRepositoryMock.On("ExistsByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(false, nil)
	userRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("domain.User")).Return(errors.New("something unexpected happened"))

	eventBusMock := new(eventmocks.Bus)

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, userApiKey, userCreatedAt)

	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_CreateUser_EventsBusError(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
  userApiKey := "test-api-key"
  userCreatedAt := "2023-10-01T00:00:00Z"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Exists", mock.Anything, mock.AnythingOfType("domain.UserID")).Return(false, nil)
  userRepositoryMock.On("ExistsByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(false, nil)
	userRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(errors.New("something unexpected happened"))

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, userApiKey, userCreatedAt)

	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_CreateUser_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
  userApiKey := "test-api-key"
  userCreatedAt := "2023-10-01T00:00:00Z"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Exists", mock.Anything, mock.AnythingOfType("domain.UserID")).Return(false, nil)
  userRepositoryMock.On("ExistsByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(false, nil)
	userRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("domain.User")).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.MatchedBy(func(events []event.Event) bool {
		evt := events[0].(usersdomain.UserCreatedEvent)
		return evt.UserName() == userName
	})).Return(nil)

	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(nil)

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, userApiKey, userCreatedAt)

	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.NoError(t, err)
}

func Test_UserService_CreateUser_UserAlreadyExists(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
  usersdomain.NewUserID(userID)
  usersdomain.NewUserName(userName)

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Exists", mock.Anything, mock.AnythingOfType("domain.UserID")).Return(true, nil)
  userRepositoryMock.On("ExistsByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(false, nil)

	eventBusMock := new(eventmocks.Bus)

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, "test-api-key", "2023-10-01T00:00:00Z")

	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, err, usersdomain.ErrUserAlreadyExists)
}

func Test_UserService_CreateUser_UserAlreadyExistsByName(t *testing.T) {
  userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
  userName := "Test User"
  usersdomain.NewUserID(userID)
  usersdomain.NewUserName(userName)

  userRepositoryMock := new(storagemocks.UserRepository)
  userRepositoryMock.On("Exists", mock.Anything, mock.AnythingOfType("domain.UserID")).Return(false, nil)
  userRepositoryMock.On("ExistsByName", mock.Anything, mock.AnythingOfType("domain.UserName")).Return(true, nil)

  eventBusMock := new(eventmocks.Bus)

  userService := NewUserService(userRepositoryMock, eventBusMock)

  err := userService.CreateUser(context.Background(), userID, userName, "test-api-key", "2023-10-01T00:00:00Z")

  userRepositoryMock.AssertExpectations(t)
  eventBusMock.AssertExpectations(t)
  assert.Error(t, err)
  assert.Equal(t, err, usersdomain.ErrUserAlreadyExists)
}