package domain

import (
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

const ExtractionCreatedEventType event.Type = "events.user.created"

type ExtractionCreatedEvent struct {
	event.BaseEvent
	id        string
	userId    string
	data      string
	metadata  string
	createdAt string
}

func NewExtractionCreatedEvent(id, userId, data, metadata, createdAt string) ExtractionCreatedEvent {
	return ExtractionCreatedEvent{
		id:        id,
		userId:    userId,
		data:      data,
		metadata:  metadata,
		createdAt: createdAt,

		BaseEvent: event.NewBaseEvent(id),
	}
}

func (e ExtractionCreatedEvent) Type() event.Type {
	return ExtractionCreatedEventType
}

func (e ExtractionCreatedEvent) ExtractionID() string {
	return e.id
}

func (e ExtractionCreatedEvent) ExtractionUserID() string {
	return e.userId
}

func (e ExtractionCreatedEvent) ExtractionData() string {
	return e.data
}

func (e ExtractionCreatedEvent) ExtractionMetadata() string {
	return e.metadata
}

func (e ExtractionCreatedEvent) ExtractionCreatedAt() string {
	return e.createdAt
}
