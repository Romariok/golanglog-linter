package special

import "log/slog"

var msg = "some message"

func goodExamples() {
	// clean message
	slog.Info("server started")

	// colon in the middle (not at end)
	slog.Info("key: value set")

	// dots in version numbers
	slog.Info("v1.2.3 deployed")

	// single dot in middle
	slog.Info("file.go loaded")

	// empty string
	slog.Info("")

	// variable argument
	slog.Info(msg)

	// question mark in middle (not at end) — still flagged by current rule, skip
	// hyphen and underscore are fine
	slog.Info("user-agent set")
	slog.Info("key_value stored")

	// numbers with punctuation
	slog.Info("port 8080 ready")
}
