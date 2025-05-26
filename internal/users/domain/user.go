package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

var ErrInvalidUserID = errors.New("invalid User ID")

type UserID struct {
	value string
}

func NewUserID(value string) (UserID, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, fmt.Errorf("%w: %s", ErrInvalidUserID, value)
	}

	return UserID{
		value: v.String(),
	}, nil
}

func (id UserID) String() string {
	return id.value
}

var ErrInvalidUserName = errors.New("invalid User Name")
var ErrEmptyUserName = errors.New("the field User Name can not be empty")

type UserName struct {
	value string
}

func NewUserName(value string) (UserName, error) {
	if value == "" {
		return UserName{}, ErrEmptyUserName
	}

	return UserName{
		value: value,
	}, nil
}

func (name UserName) String() string {
	return name.value
}

var ErrUserAlreadyExists = errors.New("user already exists")

type UserApiKey struct {
	value string
}

func NewUserApiKey(value string) (UserApiKey, error) {
	if value == "" {
		return UserApiKey{}, errors.New("the field User ApiKey can not be empty")
	}

	return UserApiKey{
		value: value,
	}, nil
}

func (key UserApiKey) String() string {
	return key.value
}

type UserCreatedAt struct {
	value string
}

func NewUserCreatedAt(value string) (UserCreatedAt, error) {
	if value == "" {
		return UserCreatedAt{}, errors.New("the field User CreatedAt can not be empty")
	}

	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return UserCreatedAt{}, errors.New("the field User CreatedAt must be a valid date in RFC3339 format (e.g. 2006-01-02T15:04:05Z07:00)")
	}

	return UserCreatedAt{
		value: value,
	}, nil
}

func (createdAt UserCreatedAt) String() string {
	return createdAt.value
}

type User struct {
	Id        UserID
	Name      UserName
	ApiKey    UserApiKey
	CreatedAt UserCreatedAt

	events []event.Event
}

type UserRepository interface {
	Save(ctx context.Context, user User) error
	Exists(ctx context.Context, id UserID) (bool, error)
	Get(ctx context.Context, id UserID) (*User, error)
	GetByName(ctx context.Context, name UserName) (*User, error)
	GetByApiKey(ctx context.Context, apiKey UserApiKey) (*User, error)
	ExistsByName(ctx context.Context, name UserName) (bool, error)
}

//mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=UserRepository

func NewUser(id, name, apiKey, createdAt string) (User, error) {
	idVO, err := NewUserID(id)
	if err != nil {
		return User{}, err
	}

	nameVO, err := NewUserName(name)
	if err != nil {
		return User{}, err
	}

	apiKeyVO, err := NewUserApiKey(apiKey)
	if err != nil {
		return User{}, err
	}

	createdAtVO, err := NewUserCreatedAt(createdAt)
	if err != nil {
		return User{}, err
	}

	user := User{
		Id:        idVO,
		Name:      nameVO,
		ApiKey:    apiKeyVO,
		CreatedAt: createdAtVO,
	}

	user.Record(NewUserCreatedEvent(idVO.String(), nameVO.String(), apiKeyVO.String(), createdAtVO.String()))
	return user, nil
}

func (u *User) SetApiKey(apiKey string) error {
  apiKeyVO, err := NewUserApiKey(apiKey)
  if err != nil {
    return err
  }

  u.ApiKey = apiKeyVO
  return nil
}

func (c *User) Record(evt event.Event) {
	c.events = append(c.events, evt)
}

func (c *User) PullEvents() []event.Event {
	evt := c.events
	c.events = []event.Event{}

	return evt
}
