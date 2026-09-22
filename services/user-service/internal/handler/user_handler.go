package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.service.UpdateUser(uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, service.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
