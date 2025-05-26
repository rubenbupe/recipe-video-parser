package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/rubenbupe/recipe-video-parser/kit/command/commandmocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateApiKeyHandler_Success(t *testing.T) {
	bus := new(commandmocks.Bus)
	bus.On("Dispatch", mock.Anything, mock.AnythingOfType("updateapikey.UserApiKeyCommand")).Return(nil)
	handler := CreateUpdateApiKeyHandler(bus)
	input := UpdateApiKeyInput{
		Name:   "name",
		ApiKey: "apikey",
	}
	err := handler(context.Background(), input)
	assert.NoError(t, err)
}

func TestUpdateApiKeyHandler_ErrorDominio(t *testing.T) {
	bus := new(commandmocks.Bus)
	handler := CreateUpdateApiKeyHandler(bus)
	input := UpdateApiKeyInput{
		Name:   "",
		ApiKey: "apikey",
	}
	err := handler(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obligatorios")
}

func TestUpdateApiKeyHandler_ErrorCamposObligatorios(t *testing.T) {
	bus := new(commandmocks.Bus)
	handler := CreateUpdateApiKeyHandler(bus)
	input := UpdateApiKeyInput{}
	err := handler(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obligatorios")
}

func TestUpdateApiKeyHandler_ErrorInterno(t *testing.T) {
	bus := new(commandmocks.Bus)
	bus.On("Dispatch", mock.Anything, mock.AnythingOfType("updateapikey.UserApiKeyCommand")).Return(errors.New("fail"))
	handler := CreateUpdateApiKeyHandler(bus)
	input := UpdateApiKeyInput{
		Name:   "name",
		ApiKey: "apikey",
	}
	err := handler(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error interno")
}
