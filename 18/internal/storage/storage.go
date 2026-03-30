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

func (s *Storage) GetForDay(userID int, date time.Time) []model.Event {
	s.RLock()
	defer s.RUnlock()
	var res []model.Event
	y, m, d := date.Date()
	for _, e := range s.events {
		ey, em, ed := e.Date.Date()
		if e.UserID == userID && y == ey && m == em && d == ed {
			res = append(res, e)
		}
	}
	return res
}

func (s *Storage) GetForWeek(userID int, date time.Time) []model.Event {
	s.RLock()
	defer s.RUnlock()
	var res []model.Event
	yr, wk := date.ISOWeek()
	for _, e := range s.events {
		eyr, ewk := e.Date.ISOWeek()
		if e.UserID == userID && yr == eyr && wk == ewk {
			res = append(res, e)
		}
	}
	return res
}

func (s *Storage) GetForMonth(userID int, date time.Time) []model.Event {
	s.RLock()
	defer s.RUnlock()
	var res []model.Event
	y, m, _ := date.Date()
	for _, e := range s.events {
		ey, em, _ := e.Date.Date()
		if e.UserID == userID && y == ey && m == em {
			res = append(res, e)
		}
	}
	return res
}
