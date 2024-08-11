package storage

import (
	pb "Auth/genproto/users"
	"Auth/models"
	"context"
)

type IStorage interface {
	Token() ITokenStorage
	Admin() IAdminStorage
	User() IUserStorage
	Close()
}

type IUserStorage interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	GetProfile(ctx context.Context) (*pb.GetProfileResponse, error)
	GetUserByEmail(ctx context.Context, email string) (string, string, error)
	GetRole(ctx context.Context, email string) (string, error)
	GetUserByID(ctx context.Context, id *pb.Id) (string, string, error)
	UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequestU) (*pb.UserResponseU, error)
}

type IAdminStorage interface {
	GetProfile(ctx context.Context, id *pb.Id) (*pb.GetProfileResponse, error)
	GetUserByEmail(ctx context.Context, email string) (string, string, error)
	FetchUsers(ctx context.Context, req *pb.Filter) (*pb.UserResponses, error)
	UpdateProfileA(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error)
	Delete(ctx context.Context, id *pb.Id) error
}

type ITokenStorage interface {
	Store(ctx context.Context, token *models.RefreshTokenDetails) error
	Delete(ctx context.Context, email string) error
	Validate(ctx context.Context, token string) (bool, error)
}
