package inmemory

import (
	"context"

	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

// QueryBus is an in-memory implementation of the query.Bus.
type QueryBus struct {
	handlers map[query.Type]query.Handler
}

// NewQueryBus initializes a new instance of QueryBus.
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[query.Type]query.Handler),
	}
}

// Dispatch implements the query.Bus interface.
func (b *QueryBus) Ask(ctx context.Context, cmd query.Query) (interface{}, error) {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil, query.ErrHandlerNotFound
	}

	return handler.Handle(ctx, cmd)
}

// Register implements the query.Bus interface.
func (b *QueryBus) Register(cmdType query.Type, handler query.Handler) {
	b.handlers[cmdType] = handler
}
