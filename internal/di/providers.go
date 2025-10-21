package di

import (
	"gorm.io/gorm"

	"github.com/miladev95/golang-project-structure/internal/config"
)

// ProvideConfig provides the application configuration
func (c *Container) ProvideConfig(cfg *config.Config) error {
	return c.Provide(func() *config.Config { return cfg })
}

// ProvideDatabase provides the database connection
func (c *Container) ProvideDatabase(cfg *config.Config) error {
	return c.Provide(func() (*gorm.DB, error) {
		return config.NewDatabase(cfg)
	})
}
