package modules

import "go.uber.org/dig"

// Module defines the interface for a DI module
type Module interface {
	// Name returns the module name
	Name() string
	// Register registers the module's dependencies
	Register(container *dig.Container) error
}