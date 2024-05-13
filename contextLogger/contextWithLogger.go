package contextLogger

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func GinContextLog(ctx *gin.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx.Request.Context())
}

func ContextLog(ctx context.Context) *zerolog.Logger {
	if gc, ok := ctx.(*gin.Context); ok {
		return GinContextLog(gc)
	}
	return zerolog.Ctx(ctx)
}
