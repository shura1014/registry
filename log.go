package registry

import (
	"context"
	"github.com/shura1014/logger"
)

const (
	debugLevel = logger.DebugLevel

	infoLevel  = logger.InfoLevel
	warnLevel  = logger.WarnLevel
	errorLevel = logger.ErrorLevel
	text       = logger.TEXT
)

var (
	l   *logger.Logger
	ctx context.Context
)

func init() {
	l = logger.Default("registry")
	ctx = context.TODO()
}

func Info(msg any, a ...any) {
	l.DoPrint(ctx, infoLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Debug(msg any, a ...any) {
	l.DoPrint(ctx, debugLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Error(msg any, a ...any) {
	l.DoPrint(ctx, errorLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func ErrorSkip(msg any, skip int, a ...any) {
	l.DoPrint(ctx, errorLevel, msg, logger.GetFileNameAndLine(skip), a...)
}

func Warn(msg any, a ...any) {
	l.DoPrint(ctx, warnLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Text(msg any, a ...any) {
	l.DoPrint(ctx, text, msg, logger.GetFileNameAndLine(0), a...)
}

func Fatal(msg any, a ...any) {
	l.DoPrint(ctx, errorLevel, msg, logger.GetFileNameAndLine(0), a...)
}
