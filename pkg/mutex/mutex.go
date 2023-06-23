package mutex

import (
	"context"
	"errors"
	"time"

	"github.com/rs/xid"
)

// custom type to avoid collision
type contextKey int

const ContextLogger contextKey = iota

// Config holds the required details for mutex
type Config struct {
	// Scope of the resource this will be prefixed with the mutex resource
	Scope string
}

// MutexProvider implements a generic interface for cache clients
type IProvider interface {
	// Get should provide the details stored under the key provided
	Get(key string) (string, error)

	// Obtain this should reserve the resource atomically
	Obtain(key string, value string, ttl time.Duration) error

	// Release should delete the resource
	Release(key string) error

	// KeepAlive to hold the resource longer
	KeepAlive(key string, ttl time.Duration) error
}

var (
	// Errors which could be raised while acquiring the mutex
	ErrorNilClient    = errors.New("nil client not accepted")
	ErrorNilProvider  = errors.New("can not operate on nil provider")
	ErrorUnauthorized = errors.New("can not release the mutex acquired by some other request")
)

// ILogger
type ILogger interface {
	Debug(string)
	Info(string)
	Error(string)
}

// Client holds the provider through which the mutex will operate
type Client struct {
	provider IProvider
	config   Config
}

// IClient mutex client interface
type IClient interface {
	New(ctx context.Context, resource string, requestID string, duration time.Duration) IMutex
}

// NewClient will set the mutex provider which will be used for further operations
func NewClient(p IProvider, conf Config) (IClient, error) {
	if p == nil {
		return nil, ErrorNilProvider
	}

	return &Client{
		provider: p,
		config:   conf,
	}, nil
}

// Mutex holds the details required to acquire the resource lock
type Mutex struct {
	// ctx holds the application context
	ctx context.Context

	// resource on which lock has to be acquired
	resource string

	// provider will hold the interface of storage provider
	provider IProvider

	// Identifier of the acquirer. This will be used while releasing the resource
	// to ensure same acquirer is releasing the resource
	// if not provided then xid will be generated and used while acquiring the resource
	acquirerID string

	// duration for which the lock has to be acquired
	// in case lock is not release within this time, then it'll be auto released
	ttl time.Duration

	// done signal will be sent once the mutex is released
	done chan bool
}

// RetryMutex holds the Mutex struct and retry arguments
type RetryMutex struct {
	*Mutex
	retries uint8
	delay   time.Duration
}

// IMutex interface for mutex
type IMutex interface {
	Lock() error
	Unlock() error
	Retry(uint8, time.Duration) IMutex
}

// New creates a new mutex struct which can then used to acquire mutex
func (c *Client) New(ctx context.Context, resource string, acquirerId string, ttl time.Duration) IMutex {
	return &Mutex{
		ctx:        ctx,
		provider:   c.provider,
		resource:   c.getScopedKey(resource),
		acquirerID: acquirerId,
		ttl:        ttl,
		done:       make(chan bool, 1),
	}
}

// Retry will add implicit retry mechanism on mutex acquire
// in case mutex was already acquired
// this will warp the mutex in retryable mutex
func (m *Mutex) Retry(retryCount uint8, retryDelay time.Duration) IMutex {
	rm := &RetryMutex{
		Mutex:   m,
		retries: retryCount,
		delay:   retryDelay,
	}

	rm.setDefaults()

	return rm
}

// Lock will set the key with request id was value only if the key in not set
// If the key is set that means the key is already acquired
func (m *Mutex) Lock() error {
	// If AcquireID is not given then set the xid
	if m.acquirerID == "" {
		m.acquirerID = xid.New().String()
	}

	if err := m.provider.Obtain(m.resource, m.acquirerID, m.ttl); err != nil {
		return err
	}

	go keepAlive(m)

	return nil
}

// Unlock will release the acquired mutex
// - this will fetch the the data of key
//      - if key does not exist that means the mutex is already released
//      - if the key has data then check the value with the request id to ensure that
//        key was set in the same request
// - delete the key of the above conditions are true
func (m *Mutex) Unlock() error {
	// send the complete signal
	m.done <- true

	val, err := m.provider.Get(m.resource)
	if err != nil {
		return err
	}

	// check if the acquirer is same as one who is calling the release
	if val != m.acquirerID {
		return ErrorUnauthorized
	}

	if err := m.provider.Release(m.resource); err != nil {
		return err
	}

	return nil
}

// keepAlive will checks the resource in intervals derived by ttl
// and increase the expire time by resource ttl
// this will ensure mutex is released by the acquirer only
func keepAlive(m *Mutex) {
	// check 5 sec before the TTL
	period := m.ttl - (time.Second * 5)

	// if period turns out less than 0 then ignore keep alive
	if period <= 0 {
		return
	}

	log, hasLogger := m.ctx.Value(ContextLogger).(ILogger)

	// keep checking in intervals for
	for {
		timer := time.NewTimer(period)
		select {
		case <-timer.C:
			if hasLogger {
				log.Debug(m.resource + ": keeping alive")
			}

			if err := m.provider.KeepAlive(m.resource, m.ttl); err != nil {
				if hasLogger {
					log.Error(m.resource + ": could not keep alive - " + err.Error())
				}
				break
			}

		case <-m.done:
			// break wont work here
			// as it breaks from the select but next loop will still continue
			// which will lead to infinite keep alive
			// so returning here to terminate the process
			return
		}
	}
}

func (rm *RetryMutex) setDefaults() {
	// if period is less than or equals 0 then ignore retries
	if rm.delay <= 0 {
		rm.retries = 0
	}
}

// RetryMutex Lock overrides Mutex Lock and has retry logic implemented
// This actually calls Mutex Lock periodically based on the retry params defined
func (rm *RetryMutex) Lock() error {
	var err error

	// Try mutex acquire periodically until acquired or all retries are exhausted
	for i := uint8(0); i <= rm.retries; i++ {
		if err = rm.Mutex.Lock(); err == nil {
			return nil
		}

		<-time.After(rm.delay)
	}

	return err
}

// getScopedKey added scope prefix to the given key
func (c *Client) getScopedKey(key string) string {
	return c.config.Scope + key
}
