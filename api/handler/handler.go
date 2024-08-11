package handler

import (
	"Auth/pkg/logger"
	"Auth/service"
	"Auth/storage"
	"log/slog"
)

type Handler struct {
	User  *service.UserService
	Log   *slog.Logger
	Admin *service.AdminService
	Token storage.ITokenStorage
}

func NewHandler(s storage.IStorage) *Handler {
	return &Handler{
		User:  service.NewUserService(s),
		Admin: service.NewAdminService(s),
		Token: s.Token(),
		Log:   logger.NewLogger(),
	}
}
