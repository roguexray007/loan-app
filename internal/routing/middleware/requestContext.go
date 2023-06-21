package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

// Set application context in middleware
func SetRequestContext(ctx *gin.Context) {
	requestID, _ := ctx.Request.Context().Value(constants.ContextKeyRequestID).(string)

	if requestID == "" {
		requestID = xid.New().String()
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(),
			constants.ContextKeyRequestID,
			requestID))
	}

	rctx, _ := tenant.Attach(ctx.Request.Context(), tenant.New().SetID(requestID))
	ctx.Request = ctx.Request.WithContext(rctx)

	ctx.Next()
}
