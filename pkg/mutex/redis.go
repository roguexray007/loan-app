package mutex

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var (
	ErrorResourceAlreadyAcquired = errors.New("resource already acquired")
)

// Redis holds the client and perform mutex related operation on it
type Redis struct {
	client redis.UniversalClient
}

// RegisterRedisClient register given redis client for further operations
func RegisterRedisClient(client redis.UniversalClient) (IProvider, error) {
	if client == nil {
		return nil, ErrorNilClient
	}

	return &Redis{
		client: client,
	}, nil
}

// Get fetches the data for the given key
// if key does not exist then returns the error
func (r *Redis) Get(key string) (string, error) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Obtain sets the resource if its not already acquired by any one.
// It will perform this atomically
func (r *Redis) Obtain(key string, value string, ttl time.Duration) error {
	if ok, err := r.client.SetNX(key, value, ttl).Result(); !ok {
		return ErrorResourceAlreadyAcquired
	} else {
		return err
	}
}

// Unlock will delete the resource
func (r *Redis) Release(key string) error {
	return r.client.Del(key).Err()
}

// Unlock will delete the resource
func (r *Redis) KeepAlive(key string, ttl time.Duration) error {
	return r.client.Expire(key, ttl).Err()
}
