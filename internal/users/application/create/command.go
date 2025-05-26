package create

import (
	"context"
	"errors"

	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

const UserCommandType command.Type = "command.user.create"

type UserCommand struct {
	id   string
	name string
  apikey string
  createdAt string
}

func NewUserCommand(id, name, apikey, createdAt string) UserCommand {
	return UserCommand{
		id:   id,
		name: name,
    apikey: apikey,
    createdAt: createdAt,
	}
}

func (c UserCommand) Type() command.Type {
	return UserCommandType
}

type UserCommandHandler struct {
	service UserService
}

func NewUserCommandHandler(service UserService) UserCommandHandler {
	return UserCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h UserCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createUserCmd, ok := cmd.(UserCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.CreateUser(
		ctx,
		createUserCmd.id,
		createUserCmd.name,
    createUserCmd.apikey,
    createUserCmd.createdAt,
	)
}

func (h UserCommandHandler) SubscribedTo() command.Type {
	return UserCommandType
}
