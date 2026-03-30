package model

import (
	"errors"
	"strings"
	"time"
)

type Event struct {
	ID    string    `json:"id"`
	Date  time.Time `json:"date"`
	Event string    `json:"event"`
}

func New(date time.Time, content string) (Event, error) {
	if strings.TrimSpace(content) == "" {
		return Event{}, errors.New("содержание события не может быть пустым")
	}

	if date.IsZero() {
		return Event{}, errors.New("дата события должна быть указана")
	}

	return Event{
		Date:  date,
		Event: content,
	}, nil
}
