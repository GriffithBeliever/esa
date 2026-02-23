package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomon/ims/internal/domain"
	"github.com/solomon/ims/internal/repository"
)

type SendInvitationInput struct {
	InviteeEmail string `json:"invitee_email"`
	Message      string `json:"message"`
}

type RespondInvitationInput struct {
	Accept bool   `json:"accept"`
	UserID string `json:"-"`
}

type InvitationService interface {
	Send(ctx context.Context, inviterID, eventID string, input SendInvitationInput) (*domain.Invitation, error)
	ListByEvent(ctx context.Context, requesterID, eventID string) ([]*domain.Invitation, error)
	ListIncoming(ctx context.Context, userEmail string) ([]*domain.Invitation, error)
	Respond(ctx context.Context, token string, input RespondInvitationInput) error
	Delete(ctx context.Context, requesterID, eventID, invID string) error
}

type invitationService struct {
	invRepo   repository.InvitationRepository
	eventRepo repository.EventRepository
	userRepo  repository.UserRepository
}

func NewInvitationService(
	invRepo repository.InvitationRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
) InvitationService {
	return &invitationService{invRepo: invRepo, eventRepo: eventRepo, userRepo: userRepo}
}

func (s *invitationService) Send(ctx context.Context, inviterID, eventID string, input SendInvitationInput) (*domain.Invitation, error) {
	if input.InviteeEmail == "" {
		return nil, fmt.Errorf("%w: invitee_email is required", domain.ErrInvalidInput)
	}

	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if !isOrganizer(e, inviterID) {
		return nil, domain.ErrForbidden
	}

	token := make([]byte, 24)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	// Look up invitee if registered
	var inviteeID *string
	if invitee, err := s.userRepo.FindByEmail(ctx, input.InviteeEmail); err == nil {
		inviteeID = &invitee.ID
	}

	now := time.Now().UTC()
	inv := &domain.Invitation{
		ID:           uuid.NewString(),
		EventID:      eventID,
		InviterID:    inviterID,
		InviteeEmail: input.InviteeEmail,
		InviteeID:    inviteeID,
		Token:        hex.EncodeToString(token),
		Status:       domain.InvitationStatusPending,
		Message:      input.Message,
		CreatedAt:    now,
		ExpiresAt:    now.Add(7 * 24 * time.Hour),
	}

	if err := s.invRepo.Create(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

func (s *invitationService) ListByEvent(ctx context.Context, requesterID, eventID string) ([]*domain.Invitation, error) {
	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if !isOrganizer(e, requesterID) {
		return nil, domain.ErrForbidden
	}
	return s.invRepo.ListByEvent(ctx, eventID)
}

func (s *invitationService) ListIncoming(ctx context.Context, userEmail string) ([]*domain.Invitation, error) {
	invs, err := s.invRepo.ListIncoming(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	// Populate events
	for _, inv := range invs {
		e, err := s.eventRepo.FindByID(ctx, inv.EventID)
		if err == nil {
			inv.Event = e
		}
	}
	return invs, nil
}

func (s *invitationService) Respond(ctx context.Context, token string, input RespondInvitationInput) error {
	inv, err := s.invRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}
	if inv.Status != domain.InvitationStatusPending {
		return domain.ErrInvitationClosed
	}
	if time.Now().After(inv.ExpiresAt) {
		return domain.ErrExpired
	}

	now := time.Now().UTC()
	inv.RespondedAt = &now

	status := domain.InvitationStatusDeclined
	if input.Accept {
		status = domain.InvitationStatusAccepted
	}

	var inviteeID *string
	if input.UserID != "" {
		inviteeID = &input.UserID
		if input.Accept {
			_ = s.eventRepo.AddAttendee(ctx, &domain.EventAttendee{
				EventID:  inv.EventID,
				UserID:   input.UserID,
				Status:   domain.EventStatusAttending,
				JoinedAt: now,
			})
		}
	}
	inv.InviteeID = inviteeID

	return s.invRepo.UpdateStatus(ctx, inv.ID, status, inv)
}

func (s *invitationService) Delete(ctx context.Context, requesterID, eventID, invID string) error {
	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if !isOrganizer(e, requesterID) {
		return domain.ErrForbidden
	}
	return s.invRepo.Delete(ctx, invID)
}

func isOrganizer(e *domain.Event, userID string) bool {
	if e.OwnerID == userID {
		return true
	}
	for _, a := range e.Attendees {
		if a.UserID == userID && a.IsOrganizer {
			return true
		}
	}
	return false
}
