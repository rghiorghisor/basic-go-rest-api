package container

import (
	"github.com/rghiorghisor/basic-go-rest-api/server"
	"go.uber.org/dig"
)

// Container is the dependency container used by the application server.
// It delegates most of it's functionality to the underlaying dig.Container.
type Container struct {
	delegate *dig.Container
}

// New Creates a new default Container.
func New() *Container {
	created := &Container{
		delegate: dig.New(),
	}

	created.Provide(server.NewControllersWithParams)

	return created
}

// Provide tells the container how to build certain types.
// (Please see the dig.Provide documentation for more details)
func (c *Container) Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	return c.delegate.Provide(constructor, opts...)
}

// Invoke runs the given function after instantiating its dependencies.
// (Please see the dig.Invoke documentation for more details)
func (c *Container) Invoke(function interface{}, opts ...dig.InvokeOption) error {
	return c.delegate.Invoke(function, opts...)
}
