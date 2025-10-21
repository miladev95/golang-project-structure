package services

import (
	"context"

	"github.com/yourusername/yourproject/internal/models"
	"github.com/yourusername/yourproject/internal/repositories"
)

// UserService defines the business logic interface for users
type UserService interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int64) error
}

// userService implements UserService
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Add business logic here (validation, etc.)
	return s.userRepo.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	// Add business logic here
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}