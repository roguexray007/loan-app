package base

import (
	"context"
	"fmt"
	"time"

	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/pkg/mutex"
)

type Core struct {
	Mutex mutex.IClient
}

// AcquireResource acquires the resource to ensure there will only one process execution
func (c Core) AcquireResource(
	ctx context.Context,
	resource string,
	duration time.Duration) (mutex.IMutex, error) {

	taskID := ctx.Value(constants.ContextKeyRequestID).(string)

	m := c.Mutex.New(ctx, resource, taskID, duration)

	if err := m.Lock(); err != nil {
		err := fmt.Errorf("failed to acquire mutex on resource: %s", resource)
		return nil, err
	}

	return m, nil
}

// ReleaseResource released the acquired mutex
func (c Core) ReleaseResource(ctx context.Context, m mutex.IMutex) {
	if err := m.Unlock(); err != nil {
		err = fmt.Errorf("failed to release mutex")
	}
}
