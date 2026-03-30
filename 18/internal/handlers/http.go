package handlers

import (
	"calendar/internal/model"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
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
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *Handler) successResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
}

func (h *Handler) Create() http.HandlerFunc {
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

func (h *Handler) Update() http.HandlerFunc {
	type req struct {
		ID     string `json:"id"`
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
		if err := h.calendar.Update(dtoIn.ID, dtoIn.UserID, date, dtoIn.Event); err != nil {
			slog.Error("failed to update event", "error", err)
			h.errorResponse(w, "failed to update event", http.StatusServiceUnavailable)
			return
		}
		h.successResponse(w, &dtoIn)

	}
}

func (h *Handler) Delete() http.HandlerFunc {
	type req struct {
		ID string `json:"id"`
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
		if err := h.calendar.Delete(dtoIn.ID); err != nil {
			slog.Error("failed to delete event", "error", err)
			h.errorResponse(w, "failed to delete event", http.StatusServiceUnavailable)
			return

		}
		h.successResponse(w, nil)
	}

}

func (h *Handler) GetForDay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Info(r.Method)
			h.errorResponse(w, "only GET method is supported", http.StatusMethodNotAllowed)
			return
		}
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			slog.Error("failed to parse user_id", "error", err)
			h.errorResponse(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			slog.Error("failed to parse date", "error", err)
			h.errorResponse(w, "invalid date", http.StatusBadRequest)
			return
		}
		res := h.calendar.GetForDay(userID, date)
		h.successResponse(w, &res)
	}
}

func (h *Handler) GetForWeek() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Info(r.Method)
			h.errorResponse(w, "only GET method is supported", http.StatusMethodNotAllowed)
			return
		}
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			slog.Error("failed to parse user_id", "error", err)
			h.errorResponse(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			slog.Error("failed to parse date", "error", err)
			h.errorResponse(w, "invalid date", http.StatusBadRequest)
			return
		}
		res := h.calendar.GetForWeek(userID, date)
		h.successResponse(w, &res)
	}
}

func (h *Handler) GetForMonth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Info(r.Method)
			h.errorResponse(w, "only GET method is supported", http.StatusMethodNotAllowed)
			return
		}
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			slog.Error("failed to parse user_id", "error", err)
			h.errorResponse(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err != nil {
			slog.Error("failed to parse date", "error", err)
			h.errorResponse(w, "invalid date", http.StatusBadRequest)
			return

		}
		res := h.calendar.GetForMonth(userID, date)
		h.successResponse(w, &res)
	}
}
