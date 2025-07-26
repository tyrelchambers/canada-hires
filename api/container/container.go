package container

import (
	"go.uber.org/dig"
)

// Container holds the dependency injection container
type Container struct {
	container *dig.Container
}

// New creates a new DI container with all dependencies registered
func New() (*Container, error) {
	c := dig.New()

	// Register all providers
	if err := registerProviders(c); err != nil {
		return nil, err
	}

	return &Container{container: c}, nil
}

// Invoke executes the given function with dependency injection
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}

// Provide registers a provider function in the container
func (c *Container) Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	return c.container.Provide(constructor, opts...)
}