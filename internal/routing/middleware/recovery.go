package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/internal/routing/response"
	"github.com/roguexray007/loan-app/pkg/utils"
)

// Recovery : Middlewares to catch uncaught error
// This will catch all un-caught panics (if any)
// In case of panic default error response will be set
func Recovery(ctx *gin.Context) {
	defer func(ctx *gin.Context) {
		if rec := recover(); rec != nil {

			err := utils.GetError(rec)

			fmt.Println(err)

			responseStruct := response.
				NewResponse(ctx).
				SetError(err).
				SetStatusCode(http.StatusInternalServerError)

			ctx.Set(constants.Response, responseStruct)
		}
	}(ctx)

	ctx.Next()
}
