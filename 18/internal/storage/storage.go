package storage

import (
	"calendar/internal/model"
	"errors"
	"sync"
)

type Storage struct {
	sync.RWMutex
	events map[string]model.Event
}

var (
	ErrEventNotFound = errors.New("событие не найдено")
)

func New() *Storage {
	return &Storage{
		events: make(map[string]model.Event),
	}
}

func (s *Storage) Create(event model.Event) error {
	s.Lock()
	defer s.Unlock()
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Update(event model.Event) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return ErrEventNotFound
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.events[id]; !ok {
		return ErrEventNotFound
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetForday()
