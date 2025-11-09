package mappers

import (
	"strings"

	"github.com/miladev95/golang-project-structure/internal/handlers/http/dtos"
	"github.com/miladev95/golang-project-structure/internal/models"
)

// ToUserResponse converts a User model to UserResponse
func ToUserResponse(user *models.User) *dtos.UserResponse {
	return &dtos.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     strings.ReplaceAll(user.Email, "@", ""),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponses converts a slice of User models to UserResponse
func ToUserResponses(users []models.User) []dtos.UserResponse {
	responses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *ToUserResponse(&user)
	}
	return responses
}