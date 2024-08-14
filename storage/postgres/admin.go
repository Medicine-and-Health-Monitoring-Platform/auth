package postgres

import (
	pb "Auth/genproto/users"
	"Auth/storage"
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type AdminRepo struct {
	DB *sql.DB
}

func NewAdminRepo(db *sql.DB) storage.IAdminStorage {

	return &AdminRepo{DB: db}
}

// GetProfile fetches the user profile by user ID
func (a *AdminRepo) GetProfile(ctx context.Context, id *pb.Id) (*pb.GetProfileResponse, error) {
	query := `
		SELECT id, email, first_name, last_name, phone_number, role, created_at
		FROM Users
		WHERE id = $1`

	var user pb.GetProfileResponse
	err := a.DB.QueryRowContext(ctx, query, id.UserId).Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
func (a *AdminRepo) GetUserByEmail(ctx context.Context, email string) (string, string, error) {
	query := `
	SELECT
		id, password_hash
	FROM
		users
	WHERE
		email = $1 and deleted_at = 0
	`
	fmt.Println("fhasdfhasdfdasfgasydfdasgfgsafasdfas fdubu gasdfyas7dfasdfsd")
	fmt.Println(a.DB)
	row := a.DB.QueryRow(query, email)
	fmt.Println("salo,")
	var id, passwordHash string
	err := row.Scan(&id, &passwordHash)
	if err != nil {
		return "", "", err
	}

	return id, passwordHash, nil
}

// UpdateProfile updates the user profile
func (a *AdminRepo) UpdateProfileA(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	query := `
		UPDATE Users
		SET first_name = $1, last_name = $2, phone_number = $3, role = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, email, first_name, last_name, phone_number, role, created_at`

	var user pb.UserResponse
	err := a.DB.QueryRowContext(ctx, query, req.FirstName, req.LastName, req.PhoneNumber, req.Role, time.Now(), req.UserId).
		Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Role,
			&user.CreatedAt,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
func (a *AdminRepo) Delete(ctx context.Context, id *pb.Id) error {
	query := `
	update
		users
	set
		deleted_at = EXTRACT(EPOCH FROM NOW())
	where
		id = $1 and deleted_at = 0 and role <> 'admin'
	`
	rows, err := a.DB.ExecContext(ctx, query, id.UserId)
	if err != nil {
		return errors.Wrap(err, "user deletion failure")
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "rows affected failure")
	}
	if rowsAffected < 1 {
		return errors.New("user not found")
	}

	return nil
}

// FetchUsers retrieves users based on filters like role, and supports pagination.
func (a *AdminRepo) FetchUsers(ctx context.Context, req *pb.Filter) (*pb.UserResponses, error) {
	query := `SELECT id, email, first_name, last_name, phone_number, role, created_at FROM Users WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	// Filtering by role
	if req.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, req.Role)
		argIndex++
	}
	// Filtering by first name
	if req.FirstName != "" {
		query += fmt.Sprintf(" AND first_name = $%d", argIndex)
		args = append(args, req.FirstName)
		argIndex++
	}

	// Pagination
	if req.Page <= 0 {
		req.Page = 1 // Default to first page
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10 // Default limit
	}
	offset := (req.Page - 1) * limit

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := a.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*pb.UserResponse{}
	for rows.Next() {
		var user pb.UserResponse
		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.UserResponses{Users: users}, nil
}
