package di

import (
	"go.uber.org/dig"

	"github.com/miladev95/golang-project-structure/internal/config"
	"github.com/miladev95/golang-project-structure/internal/di/modules"
	"github.com/miladev95/golang-project-structure/internal/handlers/http"
)

// Container represents the dependency injection container
type Container struct {
	*dig.Container
	moduleRegistry *modules.Registry
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	return &Container{
		Container:      dig.New(),
		moduleRegistry: modules.NewRegistry(),
	}
}

// RegisterModule adds a module to the container
func (c *Container) RegisterModule(module modules.Module) *Container {
	c.moduleRegistry.Register(module)
	return c
}

// Setup initializes all dependencies
func (c *Container) Setup(cfg *config.Config) error {
	// Provide core dependencies
	if err := c.ProvideConfig(cfg); err != nil {
		return err
	}

	if err := c.ProvideDatabase(cfg); err != nil {
		return err
	}

	// Setup all registered modules
	if err := c.moduleRegistry.Setup(c.Container); err != nil {
		return err
	}

	return nil
}

// GetUserHandler resolves and returns UserHandler
func (c *Container) GetUserHandler() (*http.UserHandler, error) {
	var handler *http.UserHandler
	if err := c.Invoke(func(h *http.UserHandler) {
		handler = h
	}); err != nil {
		return nil, err
	}
	return handler, nil
}

// GetModule returns a module by name (for inspection)
func (c *Container) GetModule(name string) modules.Module {
	for _, m := range c.moduleRegistry.GetModules() {
		if m.Name() == name {
			return m
		}
	}
	return nil
}
