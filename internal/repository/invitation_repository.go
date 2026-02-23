package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/solomon/ims/internal/domain"
)

type InvitationRepository interface {
	Create(ctx context.Context, inv *domain.Invitation) error
	FindByID(ctx context.Context, id string) (*domain.Invitation, error)
	FindByToken(ctx context.Context, token string) (*domain.Invitation, error)
	ListByEvent(ctx context.Context, eventID string) ([]*domain.Invitation, error)
	ListIncoming(ctx context.Context, email string) ([]*domain.Invitation, error)
	UpdateStatus(ctx context.Context, id string, status domain.InvitationStatus, inv *domain.Invitation) error
	Delete(ctx context.Context, id string) error
}

type invitationRepo struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) InvitationRepository {
	return &invitationRepo{db: db}
}

const invitationCols = `
    i.id, i.event_id, i.inviter_id, i.invitee_email, i.invitee_id,
    i.token, i.status, i.message, i.created_at, i.expires_at, i.responded_at`

func (r *invitationRepo) Create(ctx context.Context, inv *domain.Invitation) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO invitations (id, event_id, inviter_id, invitee_email, invitee_id, token, status, message, created_at, expires_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		inv.ID, inv.EventID, inv.InviterID, inv.InviteeEmail, inv.InviteeID,
		inv.Token, inv.Status, inv.Message, inv.CreatedAt, inv.ExpiresAt,
	)
	return err
}

func (r *invitationRepo) FindByID(ctx context.Context, id string) (*domain.Invitation, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+invitationCols+` FROM invitations i WHERE i.id = ?`, id)
	return r.scanInvitation(row)
}

func (r *invitationRepo) FindByToken(ctx context.Context, token string) (*domain.Invitation, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+invitationCols+` FROM invitations i WHERE i.token = ?`, token)
	return r.scanInvitation(row)
}

func (r *invitationRepo) ListByEvent(ctx context.Context, eventID string) ([]*domain.Invitation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+invitationCols+` FROM invitations i WHERE i.event_id = ? ORDER BY i.created_at DESC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanInvitations(rows)
}

func (r *invitationRepo) ListIncoming(ctx context.Context, email string) ([]*domain.Invitation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+invitationCols+` FROM invitations i
         WHERE i.invitee_email = ? AND i.status = 'pending'
         ORDER BY i.created_at DESC`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanInvitations(rows)
}

func (r *invitationRepo) UpdateStatus(ctx context.Context, id string, status domain.InvitationStatus, inv *domain.Invitation) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE invitations SET status = ?, responded_at = ?, invitee_id = ? WHERE id = ?`,
		status, inv.RespondedAt, inv.InviteeID, id,
	)
	return err
}

func (r *invitationRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM invitations WHERE id = ?`, id)
	return err
}

func (r *invitationRepo) scanInvitation(row *sql.Row) (*domain.Invitation, error) {
	inv := &domain.Invitation{}
	err := row.Scan(
		&inv.ID, &inv.EventID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
		&inv.Token, &inv.Status, &inv.Message, &inv.CreatedAt, &inv.ExpiresAt, &inv.RespondedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return inv, err
}

func (r *invitationRepo) scanInvitations(rows *sql.Rows) ([]*domain.Invitation, error) {
	var invs []*domain.Invitation
	for rows.Next() {
		inv := &domain.Invitation{}
		if err := rows.Scan(
			&inv.ID, &inv.EventID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
			&inv.Token, &inv.Status, &inv.Message, &inv.CreatedAt, &inv.ExpiresAt, &inv.RespondedAt,
		); err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, rows.Err()
}
