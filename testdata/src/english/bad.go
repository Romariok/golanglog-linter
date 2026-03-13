package english

import "log/slog"

func badExamples() {
	// Cyrillic characters
	slog.Info("запуск сервера") // want `log message must be in English only`

	// Chinese characters
	slog.Info("服务器启动") // want `log message must be in English only`

	// Mixed latin + cyrillic
	slog.Info("server запущен") // want `log message must be in English only`

	// Latin with accents
	slog.Info("café connection") // want `log message must be in English only`

	// Arabic characters
	slog.Info("خطأ") // want `log message must be in English only`

	// Japanese characters
	slog.Info("エラー") // want `log message must be in English only`

	// Non-ASCII in concatenation
	slog.Info("prefix: " + "данные") // want `log message must be in English only`
}
