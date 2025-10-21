package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/handlers/response"
	"github.com/yourusername/yourproject/internal/models"
	"github.com/yourusername/yourproject/internal/services"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}
	response.SuccessOK(c, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorBadRequest(c, "invalid id")
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		response.ErrorNotFound(c, err.Error())
		return
	}

	response.SuccessOK(c, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.ErrorBadRequest(c, err.Error())
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessCreated(c, createdUser)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorBadRequest(c, "invalid id")
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.ErrorBadRequest(c, err.Error())
		return
	}

	user.ID = id
	if err := h.userService.UpdateUser(c.Request.Context(), &user); err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessOKWithMessage(c, user, "User updated successfully")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorBadRequest(c, "invalid id")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		response.ErrorInternalServer(c, err.Error())
		return
	}

	response.SuccessNoContent(c)
}