package lowercase

import (
	"log/slog"

	"go.uber.org/zap"
)

var msg = "some message"

func goodExamples() {
	// lowercase first letter
	slog.Info("starting server")

	// starts with digit
	slog.Info("3 retries left")

	// starts with special char (rule 3 would catch, but not rule 1)
	slog.Info("!error")

	// empty string
	slog.Info("")

	// variable argument (not a string literal)
	slog.Info(msg)

	// concatenation starting with lowercase literal
	slog.Info("starting: " + "details")

	// zap lowercase
	logger := zap.NewNop()
	logger.Info("starting zap")
	logger.Debug("debug info")

	sugar := logger.Sugar()
	sugar.Info("sugar info")
	sugar.Infof("formatted %s", "value")
}
