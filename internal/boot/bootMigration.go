package boot

import (
	"context"

	"github.com/roguexray007/loan-app/internal/provider"
)

// Migrations struct will holds the data required to initiate migrations application
// this will also manage initialization of Migrations application
type Migrations struct {
	base
}

// Init will take care of initializing the Migrations with all dependencies required
func (migration *Migrations) Init(ctx context.Context) {
	migration.base.init(ctx, []string{
		provider.Config,
		provider.Database,
	})
}
