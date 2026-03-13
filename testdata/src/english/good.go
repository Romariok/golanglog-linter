package english

import "log/slog"

var msg = "some message"

func goodExamples() {
	// ASCII only
	slog.Info("server started")

	// Numbers and ASCII punctuation
	slog.Info("port 8080 ready")

	// ASCII with various punctuation
	slog.Info("key: value set")

	// Empty string
	slog.Info("")

	// Variable argument (not checked)
	slog.Info(msg)

	// Special ASCII chars
	slog.Info("v1.2.3 deployed")

	// Multiple ASCII punctuation types
	slog.Info("user (admin) logged in")
}
