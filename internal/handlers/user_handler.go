package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mohammadhprp/passport/internal/models"
	"github.com/mohammadhprp/passport/internal/repositories"
	"github.com/mohammadhprp/passport/internal/services"
)

type UserHandler struct {
	service services.UserService
}

type createUserRequest struct {
	Email         string   `json:"email" binding:"required,email"`
	Password      string   `json:"password" binding:"required,min=8"`
	MFAFactors    []string `json:"mfa_factors"`
	Status        string   `json:"status"`
	EmailVerified bool     `json:"email_verified"`
}

type userResponse struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	MFAFactors    []string `json:"mfa_factors"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", h.CreateUser)
	router.GET("", h.ListUsers)
	router.GET("/:id", h.GetUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := services.CreateUserParams{
		Email:         req.Email,
		Password:      req.Password,
		MFAFactors:    req.MFAFactors,
		EmailVerified: req.EmailVerified,
	}

	if req.Status != "" {
		params.Status = models.UserStatus(req.Status)
	}

	user, err := h.service.CreateUser(c.Request.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPassword), errors.Is(err, services.ErrInvalidUserStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, repositories.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, err := parseQueryInt(c, "limit")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := parseQueryInt(c, "offset")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	users, err := h.service.ListUsers(c.Request.Context(), services.ListUsersFilter{Limit: limit, Offset: offset})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	responses := make([]userResponse, 0, len(users))
	for i := range users {
		responses = append(responses, toUserResponse(&users[i]))
	}

	c.JSON(http.StatusOK, gin.H{"users": responses})
}

func parseQueryInt(c *gin.Context, key string) (int, error) {
	value := c.Query(key)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func toUserResponse(user *models.User) userResponse {
	factors := user.MFAFactors
	if factors == nil {
		factors = []string{}
	}

	return userResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		MFAFactors:    factors,
		Status:        string(user.Status),
		CreatedAt:     user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     user.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
