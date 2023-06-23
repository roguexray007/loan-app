package provider

import (
	"sync"

	connector "github.com/go-redis/redis"

	"github.com/roguexray007/loan-app/internal/config"
	"github.com/roguexray007/loan-app/pkg/container"
	mutexClient "github.com/roguexray007/loan-app/pkg/mutex"
)

type mutex struct {
	sync.Mutex
}

// Mutex: Key which holds the redis instance in container
const Mutex = "mutex"

func init() {
	dep.Register(Mutex, &mutex{})
}

// Build will build the new instance of logger
// it'll use the application config to initialize the redis
func (mc *mutex) Build(c container.IContainer) (container.IDependency, error) {
	mc.Lock()
	defer mc.Unlock()

	value, _ := c.Get(Redis)

	var redisClient connector.UniversalClient
	if value != nil {
		redisClient = value.(connector.UniversalClient)
	}

	configVal, _ := c.Get(Config)
	cfg := configVal.(*config.Config).Mutex

	provider, err := mutexClient.RegisterRedisClient(redisClient)
	if err != nil {
		return nil, err
	}

	client, err := mutexClient.NewClient(provider, cfg)
	if err != nil {
		return nil, err
	}

	c.Put(Mutex, client)

	return client, nil
}
