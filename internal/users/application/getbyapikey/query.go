package getbyapikey

import (
	"context"
	"errors"

	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

const UserQueryByApiKeyType query.Type = "query.user.getbyapikey"

type UserQueryByApiKey struct {
	apiKey string
}

func NewUserQueryByApiKey(apiKey string) UserQueryByApiKey {
	return UserQueryByApiKey{
		apiKey: apiKey,
	}
}

func (c UserQueryByApiKey) Type() query.Type {
	return UserQueryByApiKeyType
}

type UserQueryByApiKeyHandler struct {
	service UserByApiKeyService
}

func NewUserQueryByApiKeyHandler(service UserByApiKeyService) UserQueryByApiKeyHandler {
	return UserQueryByApiKeyHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h UserQueryByApiKeyHandler) Handle(ctx context.Context, cmd query.Query) (interface{}, error) {
	createUserCmd, ok := cmd.(UserQueryByApiKey)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	return h.service.GetUser(
		ctx,
		createUserCmd.apiKey,
	)
}

func (h UserQueryByApiKeyHandler) SubscribedTo() query.Type {
	return UserQueryByApiKeyType
}
