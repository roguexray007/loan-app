package boot

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/provider"
	"github.com/roguexray007/loan-app/internal/routing/middleware"
)

type FunctionalTest struct {
	base
}

func (api *FunctionalTest) Init(ctx context.Context) *gin.Engine {
	api.base.init(ctx, []string{
		provider.Config,
		provider.Database,
	})

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.SetRequestContext)
	router.Use(middleware.Serialize, middleware.Recovery)

	return router
}
