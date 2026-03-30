package service

import (
	"calendar/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Create(event model.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockStore) Update(event model.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStore) GetForDay(userID int, date time.Time) []model.Event {
	args := m.Called(userID, date)
	return args.Get(0).([]model.Event)
}

func (m *MockStore) GetForWeek(userID int, date time.Time) []model.Event {
	args := m.Called(userID, date)
	return args.Get(0).([]model.Event)
}

func (m *MockStore) GetForMonth(userID int, date time.Time) []model.Event {
	args := m.Called(userID, date)
	return args.Get(0).([]model.Event)
}

func TestService_Create(t *testing.T) {
	mockStore := new(MockStore)
	srv := New(mockStore)

	userID := 1
	now := time.Now()
	content := "Meeting"

	mockStore.On("Create", mock.MatchedBy(func(e model.Event) bool {
		return e.UserID == userID && e.Event == content && e.ID != ""
	})).Return(nil)

	res, err := srv.Create(userID, now, content)

	assert.NoError(t, err)
	assert.Equal(t, content, res.Event)
	assert.NotEmpty(t, res.ID)
	mockStore.AssertExpectations(t)
}

func TestService_Update(t *testing.T) {
	mockStore := new(MockStore)
	srv := New(mockStore)

	id := "test-uuid"
	userID := 1
	now := time.Now()
	content := "Updated Meeting"

	mockStore.On("Update", model.Event{
		ID:     id,
		UserID: userID,
		Date:   now,
		Event:  content,
	}).Return(nil)

	err := srv.Update(id, userID, now, content)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestService_Delete(t *testing.T) {
	mockStore := new(MockStore)
	srv := New(mockStore)
	id := "test-uuid"

	mockStore.On("Delete", id).Return(nil)

	err := srv.Delete(id)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestService_GetForDay(t *testing.T) {
	mockStore := new(MockStore)
	srv := New(mockStore)
	userID := 1
	date := time.Now()
	expected := []model.Event{{ID: "1", UserID: userID, Event: "Work"}}

	mockStore.On("GetForDay", userID, date).Return(expected)

	res := srv.GetForDay(userID, date)

	assert.Equal(t, expected, res)
	mockStore.AssertExpectations(t)
}

func TestService_GetForWeek(t *testing.T) {
	mockStore := new(MockStore)
	srv := New(mockStore)
	userID := 1
	date := time.Now()
	expected := []model.Event{{ID: "2", UserID: userID, Event: "Gym"}}

	mockStore.On("GetForWeek", userID, date).Return(expected)

	res := srv.GetForWeek(userID, date)

	assert.Equal(t, expected, res)
	mockStore.AssertExpectations(t)
}

func TestService_GetForMonth(t *testing.T) {
	mockStore := new(MockStore)
	srv := New(mockStore)
	userID := 1
	date := time.Now()
	expected := []model.Event{{ID: "3", UserID: userID, Event: "Vacation"}}

	mockStore.On("GetForMonth", userID, date).Return(expected)

	res := srv.GetForMonth(userID, date)

	assert.Equal(t, expected, res)
	mockStore.AssertExpectations(t)
}
