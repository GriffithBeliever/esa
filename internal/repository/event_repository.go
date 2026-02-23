package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/solomon/ims/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	FindByID(ctx context.Context, id string) (*domain.Event, error)
	List(ctx context.Context, userID string, filter domain.EventFilter) ([]*domain.Event, int, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error

	AddAttendee(ctx context.Context, attendee *domain.EventAttendee) error
	UpdateAttendeeStatus(ctx context.Context, eventID, userID string, status domain.EventStatus) error
	GetAttendees(ctx context.Context, eventID string) ([]domain.EventAttendee, error)
	FindConflicts(ctx context.Context, userID string, startsAt, endsAt time.Time, excludeID string) ([]*domain.Event, error)
}

type eventRepo struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(ctx context.Context, e *domain.Event) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO events (id, owner_id, title, description, location, starts_at, ends_at, is_all_day, color, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.OwnerID, e.Title, e.Description, e.Location,
		e.StartsAt, e.EndsAt, e.IsAllDay, e.Color, e.CreatedAt, e.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO event_attendees (event_id, user_id, status, is_organizer) VALUES (?, ?, ?, TRUE)`,
		e.ID, e.OwnerID, domain.EventStatusUpcoming,
	)
	if err != nil {
		return fmt.Errorf("insert organizer attendee: %w", err)
	}

	return tx.Commit()
}

func (r *eventRepo) FindByID(ctx context.Context, id string) (*domain.Event, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT e.id, e.owner_id, e.title, e.description, e.location, e.starts_at, e.ends_at,
                e.is_all_day, e.color, e.created_at, e.updated_at,
                u.id, u.email, u.username
         FROM events e
         JOIN users u ON u.id = e.owner_id
         WHERE e.id = ?`, id)

	e := &domain.Event{Owner: &domain.User{}}
	err := row.Scan(
		&e.ID, &e.OwnerID, &e.Title, &e.Description, &e.Location,
		&e.StartsAt, &e.EndsAt, &e.IsAllDay, &e.Color, &e.CreatedAt, &e.UpdatedAt,
		&e.Owner.ID, &e.Owner.Email, &e.Owner.Username,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	e.Attendees, err = r.GetAttendees(ctx, id)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *eventRepo) List(ctx context.Context, userID string, filter domain.EventFilter) ([]*domain.Event, int, error) {
	args := []any{userID}
	where := []string{"(e.owner_id = ? OR ea.user_id = ?)"}
	args = append(args, userID)

	if filter.Query != "" {
		where = append(where, "(e.title LIKE ? OR e.description LIKE ? OR e.location LIKE ?)")
		q := "%" + filter.Query + "%"
		args = append(args, q, q, q)
	}
	if filter.From != nil {
		where = append(where, "e.starts_at >= ?")
		args = append(args, *filter.From)
	}
	if filter.To != nil {
		where = append(where, "e.ends_at <= ?")
		args = append(args, *filter.To)
	}
	if filter.Location != "" {
		where = append(where, "e.location LIKE ?")
		args = append(args, "%"+filter.Location+"%")
	}
	if filter.Status != "" {
		where = append(where, "ea.status = ? AND ea.user_id = ?")
		args = append(args, filter.Status, userID)
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf(
		`SELECT COUNT(DISTINCT e.id) FROM events e
         LEFT JOIN event_attendees ea ON ea.event_id = e.id
         WHERE %s`, whereClause)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	args = append(args, filter.Limit, offset)
	query := fmt.Sprintf(
		`SELECT DISTINCT e.id, e.owner_id, e.title, e.description, e.location,
                e.starts_at, e.ends_at, e.is_all_day, e.color, e.created_at, e.updated_at
         FROM events e
         LEFT JOIN event_attendees ea ON ea.event_id = e.id
         WHERE %s
         ORDER BY e.starts_at ASC
         LIMIT ? OFFSET ?`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		e := &domain.Event{}
		if err := rows.Scan(
			&e.ID, &e.OwnerID, &e.Title, &e.Description, &e.Location,
			&e.StartsAt, &e.EndsAt, &e.IsAllDay, &e.Color, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}

	return events, total, rows.Err()
}

func (r *eventRepo) Update(ctx context.Context, e *domain.Event) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE events SET title=?, description=?, location=?, starts_at=?, ends_at=?,
         is_all_day=?, color=?, updated_at=? WHERE id=?`,
		e.Title, e.Description, e.Location, e.StartsAt, e.EndsAt,
		e.IsAllDay, e.Color, e.UpdatedAt, e.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *eventRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *eventRepo) AddAttendee(ctx context.Context, a *domain.EventAttendee) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO event_attendees (event_id, user_id, status, is_organizer)
         VALUES (?, ?, ?, ?)`,
		a.EventID, a.UserID, a.Status, a.IsOrganizer,
	)
	return err
}

func (r *eventRepo) UpdateAttendeeStatus(ctx context.Context, eventID, userID string, status domain.EventStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE event_attendees SET status = ? WHERE event_id = ? AND user_id = ?`,
		status, eventID, userID,
	)
	return err
}

func (r *eventRepo) GetAttendees(ctx context.Context, eventID string) ([]domain.EventAttendee, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT ea.event_id, ea.user_id, ea.status, ea.is_organizer, ea.joined_at,
                u.id, u.email, u.username
         FROM event_attendees ea
         JOIN users u ON u.id = ea.user_id
         WHERE ea.event_id = ?`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []domain.EventAttendee
	for rows.Next() {
		a := domain.EventAttendee{User: &domain.User{}}
		if err := rows.Scan(
			&a.EventID, &a.UserID, &a.Status, &a.IsOrganizer, &a.JoinedAt,
			&a.User.ID, &a.User.Email, &a.User.Username,
		); err != nil {
			return nil, err
		}
		attendees = append(attendees, a)
	}
	return attendees, rows.Err()
}

func (r *eventRepo) FindConflicts(ctx context.Context, userID string, startsAt, endsAt time.Time, excludeID string) ([]*domain.Event, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT e.id, e.owner_id, e.title, e.description, e.location,
                e.starts_at, e.ends_at, e.is_all_day, e.color, e.created_at, e.updated_at
         FROM events e
         JOIN event_attendees ea ON ea.event_id = e.id
         WHERE ea.user_id = ?
           AND e.id != ?
           AND e.is_all_day = FALSE
           AND ea.status NOT IN ('declined')
           AND e.starts_at < ? AND e.ends_at > ?`,
		userID, excludeID, endsAt, startsAt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		e := &domain.Event{}
		if err := rows.Scan(
			&e.ID, &e.OwnerID, &e.Title, &e.Description, &e.Location,
			&e.StartsAt, &e.EndsAt, &e.IsAllDay, &e.Color, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
