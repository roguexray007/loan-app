package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/controllers"
	"github.com/roguexray007/loan-app/internal/routing/middleware"
	"github.com/roguexray007/loan-app/internal/routing/path"
)

var loanRoutes = Route{
	group: "/v1/loans",
	middleware: []gin.HandlerFunc{
		middleware.DatabaseConnection(),
	},
	endpoints: []path.Endpoint{
		{
			http.MethodPost,
			"",
			path.PathAuthPrivate,
			controllers.LoanService.CreateLoan,
		},
		{
			http.MethodGet,
			"",
			path.PathAuthPrivate,
			controllers.LoanService.FetchLoans,
		},
		{
			http.MethodPost,
			"/approve",
			path.PathAuthAdmin,
			controllers.LoanService.ApproveLoan,
		},
		{
			http.MethodPost,
			"/pay",
			path.PathAuthPrivate,
			controllers.LoanService.PayLoan,
		},
	},
}

func init() {
	addRoutes(loanRoutes)
}
