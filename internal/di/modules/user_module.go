package modules

import (
	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/yourusername/yourproject/internal/handlers/http"
	"github.com/yourusername/yourproject/internal/repositories"
	postgresrepo "github.com/yourusername/yourproject/internal/repositories/postgres"
	"github.com/yourusername/yourproject/internal/services"
)

// UserModule represents the user domain module
type UserModule struct{}

// NewUserModule creates a new user module
func NewUserModule() Module {
	return &UserModule{}
}

// Name returns the module name
func (m *UserModule) Name() string {
	return "user"
}

// Register registers user module dependencies
func (m *UserModule) Register(container *dig.Container) error {
	// Register repository
	if err := container.Provide(func(db *gorm.DB) repositories.UserRepository {
		return postgresrepo.NewUserRepository(db)
	}); err != nil {
		return err
	}

	// Register service
	if err := container.Provide(func(userRepo repositories.UserRepository) services.UserService {
		return services.NewUserService(userRepo)
	}); err != nil {
		return err
	}

	// Register handler
	if err := container.Provide(func(userService services.UserService) *http.UserHandler {
		return http.NewUserHandler(userService)
	}); err != nil {
		return err
	}

	return nil
}