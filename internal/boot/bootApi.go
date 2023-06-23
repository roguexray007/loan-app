package boot

import (
	"context"

	"github.com/roguexray007/loan-app/internal/provider"
)

type API struct {
	base
}

func (api *API) Init(ctx context.Context) {
	api.base.init(ctx, []string{
		provider.Config,
		provider.Database,
		provider.Redis,
		provider.Mutex,
	})

	registerDefaultHandlers()
	registerApplicationHandler()
	serve(ctx)
}
