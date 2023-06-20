package constants

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) (interface{}, error, int)
