package domain

import "time"

type EventStatus string

const (
	EventStatusUpcoming  EventStatus = "upcoming"
	EventStatusAttending EventStatus = "attending"
	EventStatusMaybe     EventStatus = "maybe"
	EventStatusDeclined  EventStatus = "declined"
)

type Event struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	IsAllDay    bool      `json:"is_all_day"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Populated on read
	Attendees []EventAttendee `json:"attendees,omitempty"`
	Owner     *User           `json:"owner,omitempty"`
}

type EventAttendee struct {
	EventID     string      `json:"event_id"`
	UserID      string      `json:"user_id"`
	Status      EventStatus `json:"status"`
	IsOrganizer bool        `json:"is_organizer"`
	JoinedAt    time.Time   `json:"joined_at"`
	User        *User       `json:"user,omitempty"`
}

type EventFilter struct {
	Query    string
	From     *time.Time
	To       *time.Time
	Location string
	Status   EventStatus
	Page     int
	Limit    int
}
