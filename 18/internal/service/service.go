package service

import (
	"calendar/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventStore interface {
	Create(event model.Event) error
	Update(event model.Event) error
	Delete(id string) error
	GetForDay(date time.Time) []model.Event
	GetForWeek(date time.Time) []model.Event
	GetForMonth(date time.Time) []model.Event
}

type Service struct {
	store EventStore
}

func New(store EventStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(date time.Time, content string) (model.Event, error) {
	event, err := model.New(date, content)
	if err != nil {
		return model.Event{}, err
	}
	event.ID = uuid.New().String()
	if err := s.store.Create(event); err != nil {
		return model.Event{}, err
	}
	return event, nil
}

func (s *Service) Update(id string, date time.Time, content string) error {
	event, err := model.New(date, content)
	if err != nil {
		return err
	}
	event.ID = id
	return s.store.Update(event)
}

func (s *Service) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return s.store.Delete(id)
}

func (s *Service) GetForDay(date time.Time) []model.Event {
	return s.store.GetForDay(date)
}

func (s *Service) GetForWeek(date time.Time) []model.Event {
	return s.store.GetForWeek(date)
}

func (s *Service) GetForMonth(date time.Time) []model.Event {
	return s.store.GetForMonth(date)
}
