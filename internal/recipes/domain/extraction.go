package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rubenbupe/recipe-video-parser/kit/event"
)

var ErrInvalidExtractionID = errors.New("invalid Extraction ID")
var ErrInvalidExtractionUserID = errors.New("invalid Extraction User ID")
var ErrExtractionAlreadyExists = errors.New("extraction already exists")

type ExtractionID struct {
	value string
}

func NewExtractionID(value string) (ExtractionID, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return ExtractionID{}, fmt.Errorf("%w: %s", ErrInvalidExtractionID, value)
	}

	return ExtractionID{
		value: v.String(),
	}, nil
}

func (id ExtractionID) String() string {
	return id.value
}

type ExtractionUserID struct {
	value string
}

func NewExtractionUserID(value string) (ExtractionUserID, error) {
	if value == "" {
		return ExtractionUserID{}, errors.New("the field Extraction User ID can not be empty")
	}

	v, err := uuid.Parse(value)
	if err != nil {
		return ExtractionUserID{}, fmt.Errorf("%w: %s", ErrInvalidExtractionUserID, value)
	}

	return ExtractionUserID{
		value: v.String(),
	}, nil
}

func (id ExtractionUserID) String() string {
	return id.value
}

type ExtractionData struct {
	value json.RawMessage
}

func NewExtractionData(value string) (ExtractionData, error) {
	if value == "" {
		return ExtractionData{}, errors.New("the field Extraction Data can not be empty")
	}

	var js json.RawMessage
	if err := json.Unmarshal([]byte(value), &js); err != nil {
		return ExtractionData{}, errors.New("the field Extraction Data must be a valid JSON string")
	}

	return ExtractionData{
		value: js,
	}, nil
}

func (data ExtractionData) String() string {
	return string(data.value)
}

type ExtractionMetadata struct {
	value json.RawMessage
}

func NewExtractionMetadata(value string) (ExtractionMetadata, error) {
	if value == "" {
		return ExtractionMetadata{}, errors.New("the field Extraction Metadata can not be empty")
	}

	var js json.RawMessage
	if err := json.Unmarshal([]byte(value), &js); err != nil {
		return ExtractionMetadata{}, errors.New("the field Extraction Metadata must be a valid JSON string")
	}

	return ExtractionMetadata{
		value: js,
	}, nil
}

func (metadata ExtractionMetadata) String() string {
	return string(metadata.value)
}

type ExtractionCreatedAt struct {
	value string
}

func NewExtractionCreatedAt(value string) (ExtractionCreatedAt, error) {
	if value == "" {
		return ExtractionCreatedAt{}, errors.New("the field Extraction CreatedAt can not be empty")
	}

	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return ExtractionCreatedAt{}, errors.New("the field Extraction CreatedAt must be a valid date in RFC3339 format (e.g. 2006-01-02T15:04:05Z07:00)")
	}

	return ExtractionCreatedAt{
		value: value,
	}, nil
}

func (createdAt ExtractionCreatedAt) String() string {
	return createdAt.value
}

type Extraction struct {
	Id        ExtractionID
	UserId    ExtractionUserID
	Data      string
	Metadata  string
	CreatedAt ExtractionCreatedAt

	events []event.Event
}

type ExtractionRepository interface {
	Save(ctx context.Context, extraction Extraction) error
	Exists(ctx context.Context, id ExtractionID) (bool, error)
	Get(ctx context.Context, id ExtractionID) (*Extraction, error)
	GetByUserID(ctx context.Context, extractionId ExtractionUserID) ([]Extraction, error)
}

//mockery --case=snake --outpkg=storagemocks --output=platform/storage/storagemocks --name=ExtractionRepository

func NewExtraction(id, userId, data, metadata, createdAt string) (Extraction, error) {
	idVO, err := NewExtractionID(id)
	if err != nil {
		return Extraction{}, err
	}

	userIdVO, err := NewExtractionUserID(userId)
	if err != nil {
		return Extraction{}, err
	}

	dataVO, err := NewExtractionData(data)
	if err != nil {
		return Extraction{}, err
	}

	metadataVO, err := NewExtractionMetadata(metadata)
	if err != nil {
		return Extraction{}, err
	}

	createdAtVO, err := NewExtractionCreatedAt(createdAt)
	if err != nil {
		return Extraction{}, err
	}

	extraction := Extraction{
		Id:        idVO,
		UserId:    userIdVO,
		Data:      dataVO.String(),
		Metadata:  metadataVO.String(),
		CreatedAt: createdAtVO,
	}

	extraction.Record(NewExtractionCreatedEvent(idVO.String(), userIdVO.String(), dataVO.String(), metadataVO.String(), createdAtVO.String()))

	return extraction, nil
}

func (c *Extraction) Record(evt event.Event) {
	c.events = append(c.events, evt)
}

func (c *Extraction) PullEvents() []event.Event {
	evt := c.events
	c.events = []event.Event{}

	return evt
}
