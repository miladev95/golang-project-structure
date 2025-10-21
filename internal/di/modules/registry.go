package modules

import "go.uber.org/dig"

// Registry manages module registration
type Registry struct {
	modules []Module
}

// NewRegistry creates a new module registry
func NewRegistry() *Registry {
	return &Registry{
		modules: make([]Module, 0),
	}
}

// Register adds a module to the registry
func (r *Registry) Register(module Module) *Registry {
	r.modules = append(r.modules, module)
	return r
}

// Setup registers all modules in the container
func (r *Registry) Setup(container *dig.Container) error {
	for _, module := range r.modules {
		if err := module.Register(container); err != nil {
			return err
		}
	}
	return nil
}

// GetModules returns all registered modules
func (r *Registry) GetModules() []Module {
	return r.modules
}