package get

import (
	"context"
	"errors"

	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

const UserQueryType query.Type = "query.user.get"

type UserQuery struct {
	name string
}

func NewUserQuery(name string) UserQuery {
	return UserQuery{
		name: name,
	}
}

func (c UserQuery) Type() query.Type {
	return UserQueryType
}

type UserQueryHandler struct {
	service UserService
}

func NewUserQueryHandler(service UserService) UserQueryHandler {
	return UserQueryHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h UserQueryHandler) Handle(ctx context.Context, cmd query.Query) (interface{}, error) {
	createUserCmd, ok := cmd.(UserQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	return h.service.GetUser(
		ctx,
		createUserCmd.name,
	)
}

func (h UserQueryHandler) SubscribedTo() query.Type {
	return UserQueryType
}
