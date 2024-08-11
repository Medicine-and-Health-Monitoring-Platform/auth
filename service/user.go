package service

import (
	pb "Auth/genproto/users"
	"Auth/pkg/logger"
	"Auth/storage"
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type UserService struct {
	pb.UnimplementedAuthServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewUserService(s storage.IStorage) *UserService {
	return &UserService{
		storage: s,
		logger:  logger.NewLogger(),
	}
}

func (r *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	r.logger.Info("Registering user")
	req.Role = "patient"
	user, err := r.storage.User().Register(ctx, req)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, errors.Wrap(err, "failed to register user")
	}
	return user, nil
}

func (r *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	r.logger.Info("User login")
	res, err := r.storage.User().Login(ctx, req)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, errors.Wrap(err, "failed to login user")
	}
	return res, nil
}

func (r *UserService) GetProfileU(ctx context.Context, req *pb.Void) (*pb.GetProfileResponse, error) {
	r.logger.Info("GetUserProfile is starting")
	res, err := r.storage.User().GetProfile(ctx)
	if err != nil {
		er := errors.Wrap(err, "failed to get profile")
		r.logger.Error(er.Error())
		return nil, er
	}

	r.logger.Info("GetUserProfile has finished")
	return res, nil
}
func (r *UserService) GetUserByEmail(ctx context.Context, email string) (string, string, error) {
	r.logger.Info("GetUserByEmail is starting")
	id, passwordHash, err := r.storage.User().GetUserByEmail(ctx, email)
	if err != nil {
		er := errors.Wrap(err, "failed to get user by email")
		r.logger.Error(er.Error())
		return "", "", er
	}

	r.logger.Info("GetUserByEmail has finished")
	return id, passwordHash, nil
}
func (r *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequestU) (*pb.UserResponseU, error) {
	r.logger.Info("UpdateUserProfile is starting")
	res, err := r.storage.User().UpdateProfile(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to update profile")
		r.logger.Error(er.Error())
		return nil, er
	}

	r.logger.Info("UpdateUserProfile has finished")
	return res, nil
}

func (r *UserService) GetRole(ctx context.Context, email string) (string, error) {
	r.logger.Info("GetUserRole is starting")
	role, err := r.storage.User().GetRole(ctx, email)
	if err != nil {
		er := errors.Wrap(err, "failed to get role")
		r.logger.Error(er.Error())
		return "", er
	}
	r.logger.Info("GetUserRole has finished")
	return role, nil
}

func (r *UserService) GetUserByID(ctx context.Context, id *pb.Id) (string, string, error) {
	r.logger.Info("GetUserById is starting")
	email, password, err := r.storage.User().GetUserByID(ctx, id)
	if err != nil {
		er := errors.Wrap(err, "failed to get role")
		r.logger.Error(er.Error())
		return "", "", er
	}
	r.logger.Info("GetUserById has finished")
	return email, password, nil
}
