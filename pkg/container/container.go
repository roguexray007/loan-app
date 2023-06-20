// Package container use provider injection to create concrete type and wire the whole application together
package container

import (
	"context"
	"strings"
	"sync"
)

type IContainer interface {
	// Gives the application environment
	GetEnv() string

	// This should only be used by container and it's sub-package
	// Get instance by code from container. Only data store handler can be retrieved from container
	Get(key string) (IDependency, bool)

	// This should only be used by container and it's sub-package
	// Put value into container with code as the key. Only data store handler is saved in container
	Put(key string, value IDependency)

	// MustResolve resolved the given list of dependencies
	// will raise panic if failed to resolve
	MustResolve(list []string)

	// Destroy will call the closure procedure of provider if defined
	Destroy(list []string)
}

// Builder interface for resolving the provider
// all the dependencies added in `dependencyMap` should follow this
type IDependencyBuilder interface {
	Build(c IContainer) (IDependency, error)
}

// IDestroyer interface to be implemented to close the provider
type IDestroyer interface {
	Destroy(c IContainer)
}

// type holds the provider instance
type IDependencyManager interface {
	Register(key string, dependency IDependencyBuilder)
	Destroy(list ...string)
	Get(key string) (IDependencyBuilder, bool)
}

type IDependency interface{}

type Container struct {
	env               string
	mode              string
	ctx               context.Context
	dependencyManager IDependencyManager
	factoryMap        map[string]interface{}
	sync.Mutex
}

var c = &Container{}

// Init will initialize the application
// this should be called at least once at the time of application initialization
// Calling it again will refresh the container with new configurations provided
func Init(ctx context.Context, env string, manager IDependencyManager) IContainer {
	c = &Container{
		ctx:               ctx,
		env:               env,
		mode:              getMode(env),
		dependencyManager: manager,
		factoryMap:        make(map[string]interface{}),
	}

	return c
}

// MustResolve preload resolved the provider name provided
func (c *Container) MustResolve(list []string) {
	for _, v := range list {
		if _, ok := c.Get(v); !ok {
			panic("failed to load " + v)
		}
	}
}

// Destroy called the closure method if defined the the provider
func (c *Container) Destroy(list []string) {
	// destroy the dependency specified by the list
	for _, v := range list {
		if builder, ok := c.dependencyManager.Get(v); ok {
			if destroyer, ok := builder.(IDestroyer); ok {
				destroyer.Destroy(c)
				delete(c.factoryMap, v)
			}
		}
	}
}

// Get will resolve the provider given the key
// in case the provider is not already resolved, It'll resolve it
func (c *Container) Get(key string) (IDependency, bool) {
	value, found := c.factoryMap[key]
	if found && value != nil {
		return value, found
	}

	if builder, ok := c.dependencyManager.Get(key); ok {
		instance, err := builder.Build(c)
		if err != nil {
			return nil, false
		}

		return instance, true
	}

	return nil, false
}

// Put will add the provider to the container
// once added the same instance will be reused when called `Get` with the same key
func (c *Container) Put(code string, value IDependency) {
	c.Lock()
	defer c.Unlock()

	c.factoryMap[code] = value
}

// GetEnv will give the environment in which the application is running
func (c *Container) GetEnv() string {
	return c.env
}

// GetMode will give the mode in which the application is running
func (c *Container) GetMode() string {
	return c.mode
}

// gives the container manager
func GetContainer() IContainer {
	return c
}

// GetEnv will give the environment in which the application is running
func GetEnv() string {
	return c.GetEnv()
}

// GetMode will give the mode in which the application is running
func GetMode() string {
	return c.GetMode()
}

// Get will resolve the requested provider
func Get(key string) (interface{}, bool) {
	return c.Get(key)
}

// Set sets the provider by associating it with key
func Set(key string, instance interface{}) {
	c.Put(key, instance)
}

func getMode(env string) string {
	tokens := strings.Split(env, "-")
	if len(tokens) == 2 {
		return tokens[1]
	}

	return "dev"
}
