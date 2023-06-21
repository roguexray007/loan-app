package router

import (
	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/internal/routing/middleware"
	"github.com/roguexray007/loan-app/internal/routing/path"
	"github.com/roguexray007/loan-app/internal/routing/response"
	"github.com/roguexray007/loan-app/pkg/utils"
)

// Initialize will create the routeGroup and initialize the routes for the application
func Initialize() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	//initialize Default Routes
	initializeDefaultRoutes(router)

	return router
}

// initializeDefaultRoutes will initializes default routes required for app and workers both
func initializeDefaultRoutes(router *gin.Engine) {

	// This will register the default middleware required by the application
	// there middleware will be added to all the requests
	router.Use(middleware.SetRequestContext)
	router.Use(middleware.Serialize, middleware.Recovery)
}

// InitializeApplicationRoutes will initializes routes of all groups
// and adds middleware to the respective group as specified in the middleware group map
func InitializeApplicationRoutes(router *gin.Engine) {
	for _, routeGroup := range routeList {
		registerRouteGroup(router, routeGroup)
	}
}

// registerRouteGroup will register the group along with the route paths
// it also added the middleware to the group specified by the middleware group map
func registerRouteGroup(router *gin.Engine, routeGroup Route) {
	var middlewareList []gin.HandlerFunc

	if !utils.IsEmpty(routeGroup.middleware) {
		middlewareList = routeGroup.middleware
	}

	group := router.Group(routeGroup.group, middlewareList...)

	registerRoutes(group, routeGroup.endpoints)
}

// registerRoutes will register the routes given to the route group
func registerRoutes(router *gin.RouterGroup, endpoints []path.Endpoint) {

	for _, endpoint := range endpoints {
		router.Handle(endpoint.Method, endpoint.Path, RequestResponseHandler(endpoint))
	}
}

// decorate will decorate the request handler
// this will fetch the request input and call the handler
// It'll accept 3 parameter in the response
// 1. Response to be sent back
// 2. Error id there any
// 3. Http status code to be used in case of successful execution
// based in the parameter received response struct will be constructed
// If error is set then error will take more priority over response
func RequestResponseHandler(endpoint path.Endpoint) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if err := endpoint.Authenticate(ctx); err != nil {
			ctx.String(401, err.Error())
			return
		}

		result, ierr, statusCode := endpoint.Handler(ctx)

		responseStruct := response.NewResponse(ctx).
			SetResponse(result).
			SetError(ierr).
			SetStatusCode(statusCode)

		ctx.Set(constants.Response, responseStruct)
	}
}
