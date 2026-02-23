package domain

import "time"

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusExpired  InvitationStatus = "expired"
)

type Invitation struct {
	ID           string           `json:"id"`
	EventID      string           `json:"event_id"`
	InviterID    string           `json:"inviter_id"`
	InviteeEmail string           `json:"invitee_email"`
	InviteeID    *string          `json:"invitee_id,omitempty"`
	Token        string           `json:"token,omitempty"`
	Status       InvitationStatus `json:"status"`
	Message      string           `json:"message"`
	CreatedAt    time.Time        `json:"created_at"`
	ExpiresAt    time.Time        `json:"expires_at"`
	RespondedAt  *time.Time       `json:"responded_at,omitempty"`

	// Populated on read
	Event   *Event `json:"event,omitempty"`
	Inviter *User  `json:"inviter,omitempty"`
}
