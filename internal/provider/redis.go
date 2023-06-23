package provider

import (
	"net"
	"strconv"
	"sync"

	connector "github.com/go-redis/redis"

	"github.com/roguexray007/loan-app/internal/config"
	"github.com/roguexray007/loan-app/pkg/container"
)

type redis struct {
	sync.Mutex
}

// Redis: Key which holds the redis instance in container
const (
	Redis       = "redis"
	ClusterMode = "cluster"
)

func init() {
	dep.Register(Redis, &redis{})
}

// Build will build the new instance of logger
// it'll use the application config to initialize the redis
func (rc *redis) Build(c container.IContainer) (container.IDependency, error) {
	rc.Lock()
	defer rc.Unlock()

	var client connector.UniversalClient

	value, _ := c.Get(Config)
	cfg := value.(*config.Config).Redis

	address := net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port)))

	switch cfg.Mode {

	case ClusterMode:
		options := &connector.ClusterOptions{
			Addrs:    []string{address},
			Password: cfg.Password,
		}

		client = connector.NewClusterClient(options)

	default:
		options := &connector.Options{
			Addr:     address,
			Password: cfg.Password,
			DB:       int(cfg.Database),
		}

		client = connector.NewClient(options)
	}

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	c.Put(Redis, client)

	return client, nil
}
