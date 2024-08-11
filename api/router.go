package api

import (
	_ "Auth/api/docs"
	"Auth/api/handler"
	"Auth/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Authorazation
// @version 1.0
// @description Authorazation API
// @host localhost:8081
// @BasePath /auth
func NewRouter(s storage.IStorage) *gin.Engine {
	h := handler.NewHandler(s)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh-token", h.Refresh)
	auth.POST("/logout", h.Logout)

	return router
}
