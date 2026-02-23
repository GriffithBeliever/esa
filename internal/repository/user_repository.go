package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/solomon/ims/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)

	// Refresh tokens
	StoreRefreshToken(ctx context.Context, id, userID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (id, userID string, expiresAt time.Time, revoked bool, err error)
	RevokeRefreshToken(ctx context.Context, id string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID, u.Email, u.Username, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.scanUser(r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE id = ?`, id))
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.scanUser(r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE email = ?`, email))
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.scanUser(r.db.QueryRowContext(ctx,
		`SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE username = ?`, username))
}

func (r *userRepo) scanUser(row *sql.Row) (*domain.User, error) {
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) StoreRefreshToken(ctx context.Context, id, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) VALUES (?, ?, ?, ?)`,
		id, userID, tokenHash, expiresAt,
	)
	return err
}

func (r *userRepo) FindRefreshToken(ctx context.Context, tokenHash string) (id, userID string, expiresAt time.Time, revoked bool, err error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, expires_at, revoked FROM refresh_tokens WHERE token_hash = ?`, tokenHash)
	err = row.Scan(&id, &userID, &expiresAt, &revoked)
	if errors.Is(err, sql.ErrNoRows) {
		err = domain.ErrNotFound
	}
	return
}

func (r *userRepo) RevokeRefreshToken(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = TRUE WHERE id = ?`, id)
	return err
}

func (r *userRepo) RevokeAllUserTokens(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = ?`, userID)
	return err
}
