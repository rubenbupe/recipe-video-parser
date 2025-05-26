package domain

import (
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

const UserCreatedEventType event.Type = "events.user.created"

type UserCreatedEvent struct {
	event.BaseEvent
	id   string
	name string
  apikey string
  createdAt string
}

func NewUserCreatedEvent(id, name, apikey, createdAt string) UserCreatedEvent {
	return UserCreatedEvent{
		id:   id,
		name: name,
    apikey: apikey,
    createdAt: createdAt,

		BaseEvent: event.NewBaseEvent(id),
	}
}

func (e UserCreatedEvent) Type() event.Type {
	return UserCreatedEventType
}

func (e UserCreatedEvent) UserID() string {
	return e.id
}

func (e UserCreatedEvent) UserName() string {
	return e.name
}

func (e UserCreatedEvent) UserApiKey() string {
  return e.apikey
}

func (e UserCreatedEvent) UserCreatedAt() string {
  return e.createdAt
}
