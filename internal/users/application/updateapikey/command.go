package updateapikey

import (
	"context"
	"errors"

	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

const UserApiKeyCommandType command.Type = "command.user.updateapikey"

type UserApiKeyCommand struct {
	name   string
  apikey string
}

func NewApiKeyUserCommand(name, apikey string) UserApiKeyCommand {
	return UserApiKeyCommand{
		name: name,
    apikey: apikey,
	}
}

func (c UserApiKeyCommand) Type() command.Type {
	return UserApiKeyCommandType
}

type UserApiKeyCommandHandler struct {
	service UserApiKeyService
}

func NewUserApiKeyCommandHandler(service UserApiKeyService) UserApiKeyCommandHandler {
	return UserApiKeyCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h UserApiKeyCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createUserCmd, ok := cmd.(UserApiKeyCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.UpdateApiKey(
		ctx,
		createUserCmd.name,
    createUserCmd.apikey,
	)
}

func (h UserApiKeyCommandHandler) SubscribedTo() command.Type {
	return UserApiKeyCommandType
}
