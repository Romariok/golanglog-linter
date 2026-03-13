package sensitive

import (
	"log/slog"

	"go.uber.org/zap"
)

// Request simulates a struct with a sensitive field name.
type Request struct {
	Password string
}

func badExamples(password, token string, req Request) {
	var apiKey string

	// string literal contains "password" keyword
	slog.Info("user password: " + password) // want `log message may contain sensitive data`

	// string literal contains "token" keyword
	slog.Info("token: " + token) // want `log message may contain sensitive data`

	// string literal contains "secret" keyword (case-insensitive via regex)
	slog.Info("SECRET=" + token) // want `log message may contain sensitive data`

	// string literal contains "credential" keyword
	slog.Info("credential stored") // want `log message may contain sensitive data`

	// variable name "apiKey" matches keyword "apikey"
	slog.Debug("debug info", slog.Any("k", apiKey)) // want `log message may contain sensitive data`

	// variable name "password" matches keyword "password"
	slog.Debug("debug info", slog.Any("p", password)) // want `log message may contain sensitive data`

	// struct field "Password" matches keyword "password"
	slog.Info("request info", slog.Any("r", req.Password)) // want `log message may contain sensitive data`

	// zap: nested string literal "token" matches keyword
	logger := zap.NewNop()
	logger.Debug("debug info", zap.String("token", token)) // want `log message may contain sensitive data`

	// string literal contains "private_key"
	slog.Info("private_key loaded") // want `log message may contain sensitive data`

	// string literal contains "api_key"
	slog.Info("api_key=" + apiKey) // want `log message may contain sensitive data`
}
