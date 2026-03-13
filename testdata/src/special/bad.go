package special

import "log/slog"

func badExamples() {
	// emoji at end
	slog.Info("server started 🚀") // want `log message must not contain special characters or emojis`

	// emoji at beginning
	slog.Info("🔥 hot reload") // want `log message must not contain special characters or emojis`

	// multiple exclamation marks
	slog.Error("connection failed!!!") // want `log message must not contain special characters or emojis`

	// trailing ellipsis
	slog.Warn("something went wrong...") // want `log message must not contain special characters or emojis`

	// single trailing exclamation mark
	slog.Error("connection failed!") // want `log message must not contain special characters or emojis`

	// trailing question mark
	slog.Warn("retry failed?") // want `log message must not contain special characters or emojis`

	// trailing colon
	slog.Info("error:") // want `log message must not contain special characters or emojis`

	// newline in string
	slog.Info("line1\nline2") // want `log message must not contain special characters or emojis`

	// multiple emojis
	slog.Info("done ✅🎉") // want `log message must not contain special characters or emojis`

	// exclamation mark in middle
	slog.Info("hello! world") // want `log message must not contain special characters or emojis`

	// repeated dots
	slog.Info("loading..") // want `log message must not contain special characters or emojis`
}
