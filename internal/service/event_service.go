package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/solomon/ims/internal/domain"
	"github.com/solomon/ims/internal/repository"
)

type CreateEventInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	IsAllDay    bool      `json:"is_all_day"`
	Color       string    `json:"color"`
}

type UpdateEventInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	IsAllDay    *bool      `json:"is_all_day"`
	Color       *string    `json:"color"`
}

type EventService interface {
	Create(ctx context.Context, userID string, input CreateEventInput) (*domain.Event, error)
	Get(ctx context.Context, userID, eventID string) (*domain.Event, error)
	List(ctx context.Context, userID string, filter domain.EventFilter) ([]*domain.Event, int, error)
	Update(ctx context.Context, userID, eventID string, input UpdateEventInput) (*domain.Event, error)
	Delete(ctx context.Context, userID, eventID string) error
	UpdateStatus(ctx context.Context, userID, eventID string, status domain.EventStatus) error
	FindConflicts(ctx context.Context, userID string, startsAt, endsAt time.Time, excludeID string) ([]*domain.Event, error)
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) Create(ctx context.Context, userID string, input CreateEventInput) (*domain.Event, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("%w: title is required", domain.ErrInvalidInput)
	}
	if input.StartsAt.IsZero() || input.EndsAt.IsZero() {
		return nil, fmt.Errorf("%w: starts_at and ends_at are required", domain.ErrInvalidInput)
	}
	if input.EndsAt.Before(input.StartsAt) {
		return nil, fmt.Errorf("%w: ends_at must be after starts_at", domain.ErrInvalidInput)
	}
	if input.Color == "" {
		input.Color = "#3B82F6"
	}

	now := time.Now().UTC()
	e := &domain.Event{
		ID:          uuid.NewString(),
		OwnerID:     userID,
		Title:       input.Title,
		Description: input.Description,
		Location:    input.Location,
		StartsAt:    input.StartsAt.UTC(),
		EndsAt:      input.EndsAt.UTC(),
		IsAllDay:    input.IsAllDay,
		Color:       input.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.eventRepo.Create(ctx, e); err != nil {
		return nil, err
	}

	return s.eventRepo.FindByID(ctx, e.ID)
}

func (s *eventService) Get(ctx context.Context, userID, eventID string) (*domain.Event, error) {
	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if !s.isParticipant(e, userID) {
		return nil, domain.ErrForbidden
	}

	return e, nil
}

func (s *eventService) List(ctx context.Context, userID string, filter domain.EventFilter) ([]*domain.Event, int, error) {
	return s.eventRepo.List(ctx, userID, filter)
}

func (s *eventService) Update(ctx context.Context, userID, eventID string, input UpdateEventInput) (*domain.Event, error) {
	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if e.OwnerID != userID {
		return nil, domain.ErrForbidden
	}

	if input.Title != nil {
		e.Title = *input.Title
	}
	if input.Description != nil {
		e.Description = *input.Description
	}
	if input.Location != nil {
		e.Location = *input.Location
	}
	if input.StartsAt != nil {
		e.StartsAt = (*input.StartsAt).UTC()
	}
	if input.EndsAt != nil {
		e.EndsAt = (*input.EndsAt).UTC()
	}
	if input.IsAllDay != nil {
		e.IsAllDay = *input.IsAllDay
	}
	if input.Color != nil {
		e.Color = *input.Color
	}
	e.UpdatedAt = time.Now().UTC()

	if e.EndsAt.Before(e.StartsAt) {
		return nil, fmt.Errorf("%w: ends_at must be after starts_at", domain.ErrInvalidInput)
	}

	if err := s.eventRepo.Update(ctx, e); err != nil {
		return nil, err
	}

	return s.eventRepo.FindByID(ctx, e.ID)
}

func (s *eventService) Delete(ctx context.Context, userID, eventID string) error {
	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if e.OwnerID != userID {
		return domain.ErrForbidden
	}
	return s.eventRepo.Delete(ctx, eventID)
}

func (s *eventService) UpdateStatus(ctx context.Context, userID, eventID string, status domain.EventStatus) error {
	e, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if !s.isParticipant(e, userID) {
		return domain.ErrForbidden
	}
	return s.eventRepo.UpdateAttendeeStatus(ctx, eventID, userID, status)
}

func (s *eventService) FindConflicts(ctx context.Context, userID string, startsAt, endsAt time.Time, excludeID string) ([]*domain.Event, error) {
	return s.eventRepo.FindConflicts(ctx, userID, startsAt, endsAt, excludeID)
}

func (s *eventService) isParticipant(e *domain.Event, userID string) bool {
	if e.OwnerID == userID {
		return true
	}
	for _, a := range e.Attendees {
		if a.UserID == userID {
			return true
		}
	}
	return false
}
