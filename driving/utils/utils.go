package utils

import (
	"context"
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type mLogger struct {
	*slog.Logger
}

func (l mLogger) Printf(message string, args ...any) {
	l.Info(message, args...)
}

func RequestLog(msg string, TypeName string, Id int64, ans bool) {
	Logger.LogAttrs(context.Background(), slog.LevelInfo, msg, slog.String("TypeName", TypeName), slog.Int64("Id", Id), slog.Bool("CheckAnswer", ans))
}

func AssertOnlyOne(data []any) bool {
	for i := range data {
		if i != 0 {
			return true
		}
	}
	return false
}
