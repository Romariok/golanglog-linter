package sensitive

import "log/slog"

func goodExamples(user, msg string) {
	// neutral message — no sensitive keywords
	slog.Info("user authenticated successfully")

	// "passthrough" does not contain any sensitive keyword
	slog.Info("passthrough enabled")

	// clean server messages
	slog.Info("server started")
	slog.Info("request received")

	// variable argument — variable name "user" has no keyword
	slog.Info(msg)

	// slog structured key "name" — no sensitive keyword
	slog.Info("user info", slog.Any("name", user))

	// "password_reset" — the string contains "password" so it would actually match;
	// use a clearly safe message instead
	slog.Info("reset complete")

	// empty string
	slog.Info("")
}
