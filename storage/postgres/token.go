package postgres

import (
	"Auth/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type TokenRepo struct {
	DB *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{DB: db}
}

// Store saves a new refresh token in the database
func (t *TokenRepo) Store(ctx context.Context, token *models.RefreshTokenDetails) error {
	query := `
	insert into
		refresh_tokens (user_id, token, expires_at)
	values
		($1, $2, $3)
	`

	_, err := t.DB.ExecContext(ctx, query, token.UserID, token.Token, token.Expiry)
	if err != nil {
		return errors.Wrap(err, "refresh token storage failure")
	}

	return nil
}

// Delete removes a refresh token from the database based on the user's email
func (t *TokenRepo) Delete(ctx context.Context, email string) error {
	admin := AdminRepo{DB: t.DB}
	fmt.Println("sdfasjfasdfhasuidfoiafdasfhgasdfsdgfyasudf")
	id, _, err := admin.GetUserByEmail(ctx, email)
	fmt.Println("assalom")
	if err != nil {
		return errors.Wrap(err, "user not found")
	}

	query := `
	delete from
		refresh_tokens
	where
		user_id = $1`
	fmt.Println("salom")
	res, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "refresh token deletion failure")
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("token not found")
	}

	return nil
}

// Validate checks if a given refresh token exists in the database and is valid
func (t *TokenRepo) Validate(ctx context.Context, token string) (bool, error) {
	var n int
	query := `
	select
		1
	from
		refresh_tokens
	where
		token = $1
	`

	err := t.DB.QueryRowContext(ctx, query, token).Scan(&n)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("token not found")
		}
		return false, errors.Wrap(err, "token retrieval failure")
	}

	return true, nil
}
