package service

import (
	pb "Auth/genproto/users"
	"Auth/pkg/logger"
	"Auth/storage"
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type AdminService struct {
	pb.UnimplementedAdminServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewAdminService(s storage.IStorage) *AdminService {
	return &AdminService{
		storage: s,
		logger:  logger.NewLogger(),
	}
}

func (a *AdminService) GetProfile(ctx context.Context, UserId *pb.Id) (*pb.GetProfileResponse, error) {
	a.logger.Info("GetAdminProfile is starting")
	res, err := a.storage.Admin().GetProfile(ctx, UserId)
	if err != nil {
		er := errors.Wrap(err, "failed to get profile")
		a.logger.Error(er.Error())
		return nil, er
	}

	a.logger.Info("GetAdminProfile has finished")
	return res, nil
}
func (a *AdminService) UpdateProfileA(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	a.logger.Info("UpdateUserProfile is starting")
	res, err := a.storage.Admin().UpdateProfileA(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to update profile")
		a.logger.Error(er.Error())
		return nil, er
	}

	a.logger.Info("UpdateUserProfile has finished")
	return res, nil
}

func (a *AdminService) GetUserByEmail(ctx context.Context, email string) (string, string, error) {
	a.logger.Info("GetUserByEmail is starting")
	id, passwordHash, err := a.storage.Admin().GetUserByEmail(ctx, email)
	if err != nil {
		er := errors.Wrap(err, "failed to get user by email")
		a.logger.Error(er.Error())
		return "", "", er
	}

	a.logger.Info("GetUserByEmail has finished")
	return id, passwordHash, nil
}
func (a *AdminService) DeleteUser(ctx context.Context, req *pb.Id) (*pb.Void, error) {
	a.logger.Info("DeleteUser is starting")

	err := a.storage.Admin().Delete(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to delete user")
		a.logger.Error(er.Error())
		return nil, er
	}

	a.logger.Info("DeleteUser is finished")
	return &pb.Void{}, nil
}

func (a *AdminService) FetchUsers(ctx context.Context, req *pb.Filter) (*pb.UserResponses, error) {
	a.logger.Info("FetchUsers is starting")

	resp, err := a.storage.Admin().FetchUsers(ctx, req)
	if err != nil {
		er := errors.Wrap(err, "failed to fetch users")
		a.logger.Error(er.Error())
		return nil, er
	}

	a.logger.Info("FetchUsers is finished")
	return resp, nil
}
