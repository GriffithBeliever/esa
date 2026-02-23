package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/solomon/ims/internal/domain"
	"github.com/solomon/ims/internal/service"
)

type AIHandler struct {
	aiSvc    service.AIService
	eventSvc service.EventService
}

func NewAIHandler(aiSvc service.AIService, eventSvc service.EventService) *AIHandler {
	return &AIHandler{aiSvc: aiSvc, eventSvc: eventSvc}
}

func (h *AIHandler) GenerateDescription(w http.ResponseWriter, r *http.Request) {
	var input service.GenerateDescriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	desc, err := h.aiSvc.GenerateDescription(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AI generation failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"description": desc})
}

func (h *AIHandler) ParseEvent(w http.ResponseWriter, r *http.Request) {
	var input service.ParseEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Today == "" {
		input.Today = time.Now().UTC().Format("2006-01-02")
	}

	parsed, err := h.aiSvc.ParseEvent(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AI parsing failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, parsed)
}

func (h *AIHandler) SuggestTimes(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())

	var body struct {
		PreferredTime time.Time       `json:"preferred_time"`
		Conflicts     []*domain.Event `json:"conflicts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// If no conflicts passed, try to find them
	if len(body.Conflicts) == 0 && !body.PreferredTime.IsZero() {
		conflicts, err := h.eventSvc.FindConflicts(
			r.Context(), userID,
			body.PreferredTime,
			body.PreferredTime.Add(time.Hour),
			"",
		)
		if err == nil {
			body.Conflicts = conflicts
		}
	}

	suggestions, err := h.aiSvc.SuggestTimes(r.Context(), service.SuggestTimesInput{
		PreferredTime: body.PreferredTime,
		Conflicts:     body.Conflicts,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AI suggestion failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"suggestions": suggestions})
}
