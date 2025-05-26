package query

import (
	"context"
	"errors"
)

// Bus defines the expected behaviour from a query bus.
type Bus interface {
	// Ask is the method used to dispatch new querys.
	Ask(context.Context, Query) (interface{}, error)
	// Register is the method used to register a new query handler.
	Register(Type, Handler)
}

//mockery --case=snake --outpkg=querymocks --output=querymocks --name=Bus

// Type represents an application query type.
type Type string

// Query represents an application query.
type Query interface {
	Type() Type
}

// Handler defines the expected behaviour from a query handler.
type Handler interface {
	Handle(context.Context, Query) (interface{}, error)
	SubscribedTo() Type
}

var ErrHandlerNotFound = errors.New("handler not found")
