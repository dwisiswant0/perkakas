package log

import (
	"context"
	"github.com/kitabisa/perkakas/v3/perkakas-log/ctxkeys"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Zlogger get zerolog sublogger from context
func Zlogger(ctx context.Context) *zerolog.Logger {
	logger := &log.Logger
	if ctx.Value(ctxkeys.CtxLogger) != nil {
		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)
		logger = &l
	}

	return logger
}
