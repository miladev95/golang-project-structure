package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response envelope
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse is for paginated responses
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Message    string      `json:"message,omitempty"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int64 `json:"total_pages"`
}

// SuccessOK returns 200 OK with data
func SuccessOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessOKWithMessage returns 200 OK with data and message
func SuccessOKWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SuccessCreated returns 201 Created with data
func SuccessCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Message: "Resource created successfully",
	})
}

// SuccessCreatedWithMessage returns 201 Created with custom message
func SuccessCreatedWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// SuccessNoContent returns 204 No Content
func SuccessNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// SuccessPaginated returns 200 OK with paginated data
func SuccessPaginated(c *gin.Context, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

// ErrorBadRequest returns 400 Bad Request
func ErrorBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorUnauthorized returns 401 Unauthorized
func ErrorUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorForbidden returns 403 Forbidden
func ErrorForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorNotFound returns 404 Not Found
func ErrorNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorConflict returns 409 Conflict
func ErrorConflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorInternalServer returns 500 Internal Server Error
func ErrorInternalServer(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorUnprocessableEntity returns 422 Unprocessable Entity
func ErrorUnprocessableEntity(c *gin.Context, message string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Error:   message,
	})
}

// ErrorTooManyRequests returns 429 Too Many Requests
func ErrorTooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Error:   message,
	})
}