package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/internal/routing/response"
)

// Serialize middleware
// translates query params and request body to map which can be directly used in the controllers
// also formats the response in json format
func Serialize(ctx *gin.Context) {
	ctx.Next()
	setResponse(ctx)
}

// sets the response and status code for the current request in json format
// data will be taken form the context value set by controller
func setResponse(ctx *gin.Context) {
	responseBody, status := prepareResponse(ctx)
	ctx.JSON(status, responseBody)
}

// if not available give the empty map
func prepareResponse(ctx *gin.Context) (interface{}, int) {
	data, exists := ctx.Get(constants.Response)

	if !exists {
		return map[string]interface{}{}, http.StatusOK
	}

	responseStruct := data.(response.Response)
	responseStruct.ModifyError(ctx.Request.Context())
	setHeaders(ctx, responseStruct.Headers())
	responseStruct.Log()

	return responseStruct.GetResponseBody(), responseStruct.StatusCode()
}

// setHeaders: will set the response header
func setHeaders(ctx *gin.Context, headers map[string]string) {
	for key, val := range headers {
		ctx.Header(key, val)
	}
}
