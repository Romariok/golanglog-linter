package lowercase

import (
	"log/slog"

	"go.uber.org/zap"
)

func badExamples() {
	// uppercase ASCII first letter
	slog.Info("Starting server") // want `log message should start with a lowercase letter`

	// uppercase unicode first letter
	slog.Info("Ñew server") // want `log message should start with a lowercase letter`

	// uppercase in string concatenation (leftmost literal checked)
	slog.Info("Error: " + "some error") // want `log message should start with a lowercase letter`

	// various slog methods
	slog.Debug("Debug message") // want `log message should start with a lowercase letter`
	slog.Warn("Warning here")   // want `log message should start with a lowercase letter`
	slog.Error("Critical fail") // want `log message should start with a lowercase letter`

	// zap Logger
	logger := zap.NewNop()
	logger.Info("Starting zap") // want `log message should start with a lowercase letter`
	logger.Debug("Debug zap")   // want `log message should start with a lowercase letter`

	// zap SugaredLogger
	sugar := logger.Sugar()
	sugar.Info("SugaredInfo") // want `log message should start with a lowercase letter`
	sugar.Infof("Formatted")  // want `log message should start with a lowercase letter`
}
