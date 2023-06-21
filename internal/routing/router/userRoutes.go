package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/controllers"
	"github.com/roguexray007/loan-app/internal/routing/middleware"
	"github.com/roguexray007/loan-app/internal/routing/path"
)

var userRoutes = Route{
	group: "/v1/users",
	middleware: []gin.HandlerFunc{
		middleware.DatabaseConnection(),
	},
	endpoints: []path.Endpoint{
		{
			http.MethodPost,
			"",
			path.PathAuthPublic,
			controllers.UserService.CreateUser,
		},
	},
}

func init() {
	addRoutes(userRoutes)
}
