package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/solomon/ims/internal/domain"
	"github.com/solomon/ims/internal/repository"
	"github.com/solomon/ims/internal/service"
)

type InvitationHandler struct {
	invSvc   service.InvitationService
	userRepo repository.UserRepository
}

func NewInvitationHandler(invSvc service.InvitationService, userRepo repository.UserRepository) *InvitationHandler {
	return &InvitationHandler{invSvc: invSvc, userRepo: userRepo}
}

func (h *InvitationHandler) Send(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")

	var input service.SendInvitationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	inv, err := h.invSvc.Send(r.Context(), userID, eventID, input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "only organizers can send invitations")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to send invitation")
		return
	}

	writeJSON(w, http.StatusCreated, inv)
}

func (h *InvitationHandler) ListByEvent(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")

	invs, err := h.invSvc.ListByEvent(r.Context(), userID, eventID)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to list invitations")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"invitations": invs})
}

func (h *InvitationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	eventID := chi.URLParam(r, "id")
	invID := chi.URLParam(r, "invID")

	if err := h.invSvc.Delete(r.Context(), userID, eventID, invID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete invitation")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InvitationHandler) ListIncoming(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromCtx(r.Context())
	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	invs, err := h.invSvc.ListIncoming(r.Context(), user.Email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list invitations")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"invitations": invs})
}

func (h *InvitationHandler) Respond(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	var body struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Optionally get userID if authenticated
	userID := UserIDFromCtx(r.Context())

	input := service.RespondInvitationInput{
		Accept: body.Accept,
		UserID: userID,
	}

	if err := h.invSvc.Respond(r.Context(), token, input); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "invitation not found")
			return
		}
		if errors.Is(err, domain.ErrExpired) {
			writeError(w, http.StatusGone, "invitation has expired")
			return
		}
		if errors.Is(err, domain.ErrInvitationClosed) {
			writeError(w, http.StatusConflict, "invitation already responded to")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to respond to invitation")
		return
	}

	status := "declined"
	if body.Accept {
		status = "accepted"
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": status})
}
