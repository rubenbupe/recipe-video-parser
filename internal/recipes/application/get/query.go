package get

import (
	"context"
	"errors"

	"github.com/rubenbupe/recipe-video-parser/kit/query"
)

const ExtractionQueryType query.Type = "query.extraction.get"

type ExtractionQuery struct {
	userId string
}

func NewExtractionQuery(userId string) ExtractionQuery {
	return ExtractionQuery{
		userId: userId,
	}
}

func (c ExtractionQuery) Type() query.Type {
	return ExtractionQueryType
}

type ExtractionQueryHandler struct {
	service ExtractionService
}

func NewExtractionQueryHandler(service ExtractionService) ExtractionQueryHandler {
	return ExtractionQueryHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h ExtractionQueryHandler) Handle(ctx context.Context, cmd query.Query) (interface{}, error) {
	getExtractionCmd, ok := cmd.(ExtractionQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	return h.service.GetExtraction(
		ctx,
		getExtractionCmd.userId,
	)
}

func (h ExtractionQueryHandler) SubscribedTo() query.Type {
	return ExtractionQueryType
}
