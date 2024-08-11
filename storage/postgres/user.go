package postgres

import (
	pb "Auth/genproto/users"
	"Auth/storage"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) storage.IUserStorage {
	return &UserRepo{DB: db}
}

// Register a new user
func (u *UserRepo) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	query := `
		INSERT INTO Users (email, password_hash, first_name, last_name, phone_number, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	var user pb.UserResponse
	err := u.DB.QueryRowContext(ctx, query, req.Email, req.Password, req.FirstName, req.LastName, req.PhoneNumber, req.Role).
		Scan(&user.Id, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.PhoneNumber = req.PhoneNumber
	user.Role = req.Role

	return &user, nil
}

// Login a user
func (u *UserRepo) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	query := `SELECT id, password_hash, role FROM Users WHERE email = $1`

	var (
		userId         string
		hashedPassword string
		role           string
	)
	err := u.DB.QueryRowContext(ctx, query, req.Email).Scan(&userId, &hashedPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate access and refresh tokens (implementation skipped for brevity)
	accessToken := "access_token_example"
	refreshToken := "refresh_token_example"

	return &pb.LoginResponse{
		Access:  accessToken,
		Refresh: refreshToken,
		UserId:  userId,
		Role:    role,
	}, nil
}

// GetProfile fetches the user profile by user ID
func (u *UserRepo) GetProfile(ctx context.Context) (*pb.GetProfileResponse, error) {
	query := `
		SELECT id, email, first_name, last_name, phone_number, role, created_at
		FROM Users
		WHERE id = $1`
	Id := ctx.Value("user_id").(string)
	var user pb.GetProfileResponse
	err := u.DB.QueryRowContext(ctx, query, Id).Scan(
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
func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (string, string, error) {
	query := `
	SELECT
		id, password_hash
	FROM
		users
	WHERE
		email = $1 and deleted_at = 0
	`
	row := u.DB.QueryRowContext(ctx, query, email)

	var id, passwordHash string
	err := row.Scan(&id, &passwordHash)
	if err != nil {
		return "", "", err
	}

	return id, passwordHash, nil
}

// UpdateProfile updates the user profile
func (u *UserRepo) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequestU) (*pb.UserResponseU, error) {
	query := `
		UPDATE Users
		SET first_name = $1, last_name = $2, phone_number = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, email, first_name, last_name, phone_number, created_at`

	var user pb.UserResponseU
	err := u.DB.QueryRowContext(ctx, query, req.FirstName, req.LastName, req.PhoneNumber, time.Now(), req.UserId).
		Scan(
			&user.Id,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
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
func (u *UserRepo) GetRole(ctx context.Context, email string) (string, error) {
	query := `
	select
		role
	from
		users
	where
		email = $1 and deleted_at = 0
	`

	var role string
	err := u.DB.QueryRowContext(ctx, query, email).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}
func (u *UserRepo) GetUserByID(ctx context.Context, id *pb.Id) (string, string, error) {
	query := `
	SELECT
		email, password_hash
	FROM
		users
	WHERE
		id = $1 and deleted_at = 0`
	row := u.DB.QueryRowContext(ctx, query, id)

	var username, email, passwordHash string
	err := row.Scan(&username, &email, &passwordHash)
	if err != nil {
		return "", "", err
	}

	return email, passwordHash, nil
}
