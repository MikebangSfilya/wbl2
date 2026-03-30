package model

import (
	"errors"
	"strings"
	"time"
)

type Event struct {
	ID     string    `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Event  string    `json:"event"`
}

func New(userID int, date time.Time, content string) (Event, error) {
	if userID <= 0 {
		return Event{}, errors.New("user_id должен быть положительным числом")
	}
	if strings.TrimSpace(content) == "" {
		return Event{}, errors.New("содержание события не может быть пустым")
	}
	if date.IsZero() {
		return Event{}, errors.New("дата события должна быть указана")
	}

	return Event{
		UserID: userID,
		Date:   date,
		Event:  content,
	}, nil
}
