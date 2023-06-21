package router

import (
	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/routing/path"
)

var routeList []Route

type Route struct {
	group      string
	middleware []gin.HandlerFunc
	endpoints  []path.Endpoint
}

func addRoutes(routes Route) {
	routeList = append(routeList, routes)
}
