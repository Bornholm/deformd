package module

import (
	"context"
	"fmt"

	"gitlab.com/wpetit/goweb/logger"
)

type LogFunc func(ctx context.Context, message string, args ...any)

func defaultLogger(ctx context.Context, message string, args ...any) {
	logger.Debug(ctx, fmt.Sprintf(message, args...))
}
