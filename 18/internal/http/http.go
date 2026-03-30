package http

import (
	"calendar/internal/model"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Calendar interface {
	Create(userID int, date time.Time, content string) (model.Event, error)
	Update(id string, userID int, date time.Time, content string) error
	Delete(id string) error
	GetForDay(userID int, date time.Time) []model.Event
	GetForWeek(userID int, date time.Time) []model.Event
	GetForMonth(userID int, date time.Time) []model.Event
}

type Handler struct {
	calendar Calendar
}

func New(calendar Calendar) *Handler {
	return &Handler{calendar: calendar}
}

func (h *Handler) errorResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *Handler) successResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
}

func (h *Handler) ServeHTTP() http.HandlerFunc {
	type req struct {
		UserID int    `json:"user_id"`
		Date   string `json:"date"`
		Event  string `json:"event"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			slog.Info(r.Method)
			h.errorResponse(w, "only POST method is supported", http.StatusMethodNotAllowed)
			return
		}
		var dtoIn req
		if err := json.NewDecoder(r.Body).Decode(&dtoIn); err != nil {
			slog.Error("failed to decode request body", "error", err)
			h.errorResponse(w, "invalid json", http.StatusBadRequest)
			return
		}

		date, err := time.Parse("2006-01-02", dtoIn.Date)
		if err != nil {
			slog.Error("failed to parse date", "error", err)
			h.errorResponse(w, "invalid date", http.StatusBadRequest)
			return
		}
		res, err := h.calendar.Create(dtoIn.UserID, date, dtoIn.Event)
		if err != nil {
			slog.Error("failed to create event", "error", err)
			h.errorResponse(w, "failed to create event", http.StatusServiceUnavailable)
			return
		}
		h.successResponse(w, &res)
	}
}
