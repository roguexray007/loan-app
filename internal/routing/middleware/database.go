package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/pkg/db"
)

func DatabaseConnection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method

		if method == http.MethodGet {
			ctx.Request = ctx.Request.WithContext(
				context.WithValue(ctx.Request.Context(),
					db.ContextKeyDatabaseConnection,
					db.Replica))
		}

		ctx.Next()
	}
}
