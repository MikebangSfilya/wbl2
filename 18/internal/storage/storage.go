package storage

import (
	"calendar/internal/model"
	"errors"
	"sync"
	"time"
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

func (s *Storage) GetForDay(date time.Time) []model.Event {
	s.RLock()
	defer s.RUnlock()

	var result []model.Event
	y, m, d := date.Date()
	for _, event := range s.events {
		ey, em, ed := event.Date.Date()
		if y == ey && m == em && d == ed {
			result = append(result, event)
		}
	}
	return result
}

func (s *Storage) GetForWeek(date time.Time) []model.Event {
	s.RLock()
	defer s.RUnlock()

	var result []model.Event
	year, week := date.ISOWeek()
	for _, event := range s.events {
		eyear, eweek := event.Date.ISOWeek()
		if year == eyear && week == eweek {
			result = append(result, event)
		}

	}
	return result
}

func (s *Storage) GetForMonth(date time.Time) []model.Event {
	s.RLock()
	defer s.RUnlock()

	var result []model.Event
	y, m, _ := date.Date()
	for _, event := range s.events {
		eyear, emonth, _ := event.Date.Date()
		if y == eyear && m == emonth {
			result = append(result, event)
		}
	}
	return result
}
