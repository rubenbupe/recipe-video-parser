package create

import (
	"context"
	"errors"

	"github.com/rubenbupe/recipe-video-parser/kit/command"
)

const ExtractionCommandType command.Type = "command.extraction.create"

type ExtractionCommand struct {
	id        string
	userId    string
	data      string
	metadata  string
	createdAt string
}

func NewExtractionCommand(id, userId, data, metadata, createdAt string) ExtractionCommand {
	return ExtractionCommand{
		id:        id,
		userId:    userId,
		data:      data,
		metadata:  metadata,
		createdAt: createdAt,
	}
}

func (c ExtractionCommand) Type() command.Type {
	return ExtractionCommandType
}

type ExtractionCommandHandler struct {
	service ExtractionService
}

func NewExtractionCommandHandler(service ExtractionService) ExtractionCommandHandler {
	return ExtractionCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h ExtractionCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	createExtractionCmd, ok := cmd.(ExtractionCommand)
	if !ok {
		return errors.New("unexpected command")
	}

	return h.service.CreateExtraction(
		ctx,
		createExtractionCmd.id,
    createExtractionCmd.userId,
    createExtractionCmd.data,
    createExtractionCmd.metadata,
    createExtractionCmd.createdAt,
	)
}

func (h ExtractionCommandHandler) SubscribedTo() command.Type {
	return ExtractionCommandType
}
