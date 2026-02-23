package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/solomon/ims/internal/domain"
	"github.com/solomon/ims/internal/service"
)

type EventHandler struct {
	eventSvc service.EventService
}

func NewEventHandler(eventSvc service.EventService) *EventHandler {
	return &EventHandler{eventSvc: eventSvc}
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	q := r.URL.Query()

	filter := domain.EventFilter{
		Query:    q.Get("q"),
		Location: q.Get("location"),
		Status:   domain.EventStatus(q.Get("status")),
	}

	if v := q.Get("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			filter.From = &t
		}
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			filter.To = &t
		}
	}
	if v := q.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Page = n
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}

	events, total, err := h.eventSvc.List(r.Context(), userID, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list events")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"events": events,
		"total":  total,
		"page":   filter.Page,
		"limit":  filter.Limit,
	})
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())

	var input service.CreateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	e, err := h.eventSvc.Create(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	writeJSON(w, http.StatusCreated, e)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")

	e, err := h.eventSvc.Get(r.Context(), userID, eventID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get event")
		return
	}

	writeJSON(w, http.StatusOK, e)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")

	var input service.UpdateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	e, err := h.eventSvc.Update(r.Context(), userID, eventID, input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "only the owner can update this event")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update event")
		return
	}

	writeJSON(w, http.StatusOK, e)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")

	if err := h.eventSvc.Delete(r.Context(), userID, eventID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "only the owner can delete this event")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")

	var body struct {
		Status domain.EventStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.eventSvc.UpdateStatus(r.Context(), userID, eventID, body.Status); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update status")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": string(body.Status)})
}

func (h *EventHandler) CheckConflicts(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	q := r.URL.Query()

	startsAt, err1 := time.Parse(time.RFC3339, q.Get("starts_at"))
	endsAt, err2 := time.Parse(time.RFC3339, q.Get("ends_at"))
	if err1 != nil || err2 != nil {
		writeError(w, http.StatusBadRequest, "starts_at and ends_at are required (RFC3339)")
		return
	}

	excludeID := q.Get("exclude_id")
	conflicts, err := h.eventSvc.FindConflicts(r.Context(), userID, startsAt, endsAt, excludeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "conflict check failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"conflicts": conflicts})
}
