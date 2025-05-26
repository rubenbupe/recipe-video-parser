package handlers

import (
	"context"
	"errors"
	"testing"

	usersdomain "github.com/rubenbupe/recipe-video-parser/internal/users/domain"
	"github.com/rubenbupe/recipe-video-parser/kit/command/commandmocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateHandler_Success(t *testing.T) {
	bus := new(commandmocks.Bus)
	bus.On("Dispatch", mock.Anything, mock.AnythingOfType("create.UserCommand")).Return(nil)
	handler := CreateHandler(bus)
	input := CreateUserInput{
		ID:        "id",
		Name:      "name",
		ApiKey:    "apikey",
		CreatedAt: "2023-01-01T00:00:00Z",
	}
	err := handler(context.Background(), input)
	assert.NoError(t, err)
}

func TestCreateHandler_ErrorDominio(t *testing.T) {
	bus := new(commandmocks.Bus)
	bus.On("Dispatch", mock.Anything, mock.AnythingOfType("create.UserCommand")).Return(usersdomain.ErrUserAlreadyExists)
	handler := CreateHandler(bus)
	input := CreateUserInput{
		ID:        "id",
		Name:      "name",
		ApiKey:    "apikey",
		CreatedAt: "2023-01-01T00:00:00Z",
	}
	err := handler(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error de dominio")
}

func TestCreateHandler_ErrorCamposObligatorios(t *testing.T) {
	bus := new(commandmocks.Bus)
	handler := CreateHandler(bus)
	input := CreateUserInput{}
	err := handler(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obligatorios")
}

func TestCreateHandler_ErrorInterno(t *testing.T) {
	bus := new(commandmocks.Bus)
	bus.On("Dispatch", mock.Anything, mock.AnythingOfType("create.UserCommand")).Return(errors.New("fail"))
	handler := CreateHandler(bus)
	input := CreateUserInput{
		ID:        "id",
		Name:      "name",
		ApiKey:    "apikey",
		CreatedAt: "2023-01-01T00:00:00Z",
	}
	err := handler(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error interno")
}
