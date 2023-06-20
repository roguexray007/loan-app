package provider

import (
	"sync"

	"github.com/roguexray007/loan-app/pkg/container"
)

// Manager holds the provider list
// gives flexibility to register new provider to the list
// and fetch the same for resolution purpose
type Manager struct {
	dependencyMap map[string]container.IDependencyBuilder
	sync.Mutex
}

var dep = &Manager{
	dependencyMap: make(map[string]container.IDependencyBuilder),
}

// GetManager gives the provider manager instance
func GetManager() container.IDependencyManager {
	return dep
}

// Destroy wrapper method to call destroy on provider manager
func Destroy(list ...string) {
	GetManager().Destroy(list...)
}

// Register will add the new provider to the list (associated with key)
func (m *Manager) Register(key string, instance container.IDependencyBuilder) {
	m.Lock()
	defer m.Unlock()

	m.dependencyMap[key] = instance
}

// Destroy close all the provider if there is a destroy method defined
func (m *Manager) Destroy(list ...string) {
	m.Lock()
	defer m.Unlock()

	// if no key provided then collect all the registered keys
	if len(list) == 0 {
		list = make([]string, len(m.dependencyMap))

		for k := range m.dependencyMap {
			list = append(list, k)
		}
	}

	container.GetContainer().Destroy(list)
}

// Get will give the provider builder for the given key
func (m *Manager) Get(key string) (container.IDependencyBuilder, bool) {
	builder, ok := m.dependencyMap[key]

	return builder, ok
}
