package service

import (
	"calendar/internal/model"
	"time"

	"github.com/google/uuid"
)

type EventStore interface {
	Create(event model.Event) error
	Update(event model.Event) error
	Delete(id string) error
	GetForDay(userID int, date time.Time) []model.Event
	GetForWeek(userID int, date time.Time) []model.Event
	GetForMonth(userID int, date time.Time) []model.Event
}

type Service struct {
	store EventStore
}

func New(store EventStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(userID int, date time.Time, content string) (model.Event, error) {
	ev, err := model.New(userID, date, content)
	if err != nil {
		return model.Event{}, err
	}
	ev.ID = uuid.New().String()
	if err := s.store.Create(ev); err != nil {
		return model.Event{}, err
	}
	return ev, nil
}

func (s *Service) Update(id string, userID int, date time.Time, content string) error {
	ev, err := model.New(userID, date, content)
	if err != nil {
		return err
	}
	ev.ID = id
	return s.store.Update(ev)
}

func (s *Service) Delete(id string) error {
	return s.store.Delete(id)
}

func (s *Service) GetForDay(userID int, date time.Time) []model.Event {
	return s.store.GetForDay(userID, date)
}

func (s *Service) GetForWeek(userID int, date time.Time) []model.Event {
	return s.store.GetForWeek(userID, date)
}

func (s *Service) GetForMonth(userID int, date time.Time) []model.Event {
	return s.store.GetForMonth(userID, date)
}
