package boot

import (
	"context"

	"github.com/roguexray007/loan-app/internal/provider"
)

type base struct {
	ctx context.Context
}

func (b *base) init(ctx context.Context, preload []string) {
	b.ctx = ctx
	provider.MustResolve(preload)
}
