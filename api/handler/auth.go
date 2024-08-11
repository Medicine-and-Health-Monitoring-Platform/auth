package handler

import (
	"Auth/api/tokens"
	pb "Auth/genproto/users"
	"Auth/models"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Registers user
// @Description Registers a new user
// @Tags auth
// @Param user body models.RegisterRequest true "User data"
// @Success 200 {object} users.UserResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	h.Log.Info("Register function is starting")

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error("Invalid data provided", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	passByte, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error("Error hashing password", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}
	req.Password = string(passByte)

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	resp, err := h.User.Register(ctx, &pb.RegisterRequest{
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		h.Log.Error("Error registering user", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	h.Log.Info("Register has successfully finished")
	c.JSON(http.StatusOK, gin.H{"New user": resp})
}

// Login godoc
// @Summary Logs user in
// @Description Logs user in
// @Tags auth
// @Param data body models.LoginRequest true "User credentials"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	h.Log.Info("Login function is starting")

	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		er := errors.Wrap(err, "invalid data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second) // Increased timeout
	defer cancel()

	id, passwordHash, err := h.User.GetUserByEmail(ctx, req.Email)
	if err != nil {
		er := errors.Wrap(err, "user not found").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}
	println(passwordHash)
	println(req.Password)

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		er := errors.Wrap(err, "invalid password").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	role, err := h.User.GetRole(ctx, req.Email)
	if err != nil {
		er := errors.Wrap(err, "error getting user role").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	accessToken, err := tokens.GenerateAccessToken(id, req.Email, role)
	if err != nil {
		er := errors.Wrap(err, "error generating access token").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	refreshToken, err := tokens.GenerateRefreshToken(id)
	if err != nil {
		er := errors.Wrap(err, "error generating refresh token").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	exp, err := tokens.GetRefreshTokenExpiry(refreshToken)
	if err != nil {
		er := errors.Wrap(err, "error getting refresh token expiry").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	// Debugging log to inspect token details
	h.Log.Info("Storing refresh token", slog.String("user_id", id), slog.String("refresh_token", refreshToken), slog.String("expiry", exp))

	err = h.Token.Store(ctx, &models.RefreshTokenDetails{
		UserID: id,
		Token:  refreshToken,
		Expiry: exp,
	})
	if err != nil {
		er := errors.Wrap(err, "error storing refresh token").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("Login has successfully finished")
	c.JSON(http.StatusOK, gin.H{"Tokens": models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}})
}

// Refresh godoc
// @Summary Refreshes refresh token
// @Description Refreshes refresh token
// @Tags auth
// @Param data body models.RefreshToken true "Refresh token"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error while processing request"
// @Router /refresh-token [post]
func (h Handler) Refresh(c *gin.Context) {
	h.Log.Info("Refresh function is starting")

	var t models.RefreshToken
	if err := c.ShouldBind(&t); err != nil {
		er := errors.Wrap(err, "invalid data").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	valid, err := tokens.ValidateRefreshToken(t.Token)
	if !valid || err != nil {
		er := errors.Wrap(err, "invalid refresh token").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	valid, err = h.Token.Validate(ctx, t.Token)
	if !valid || err != nil {
		er := errors.Wrap(err, "invalid refresh token").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	id, err := tokens.GetUserIdFromRefreshToken(t.Token)
	if err != nil {
		er := errors.Wrap(err, "error getting user id").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	email, _, err := h.User.GetUserByID(ctx, &pb.Id{UserId: id})
	if err != nil {
		er := errors.Wrap(err, "user not found").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	role, err := h.User.GetRole(ctx, email)
	if err != nil {
		er := errors.Wrap(err, "error getting user role").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	accessToken, err := tokens.GenerateAccessToken(id, email, role)
	if err != nil {
		er := errors.Wrap(err, "error generating access token").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": er},
		)
		h.Log.Error(er)
		return
	}

	h.Log.Info("Refresh has successfully finished")
	c.JSON(http.StatusOK, gin.H{"Tokens": models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: t.Token,
	}})
}

// Logout godoc
// @Summary Logouts user
// @Description Logouts user by ID
// @Tags auth
// @Param email query string true "User email"
// @Success 200 {string} string "User logged out successfully"
// @Failure 400 {object} string "Invalid user id"
// @Failure 500 {object} string "Server error while processing request"
// @Router /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	h.Log.Info("Logout function is starting")

	email := c.Query("email")
	if email == "" {
		er := errors.New("invalid email").Error()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()

	err := h.Token.Delete(ctx, email)
	fmt.Println("sadfjadfdasfads gadsyf asdfa d8uf9udfugyasdf")
	if err != nil {
		er := errors.Wrap(err, "error logging out").Error()
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": er})
		h.Log.Error(er)
		return
	}

	h.Log.Info("Logout has successfully finished")
	c.JSON(http.StatusOK, "User logged out successfully")

}
