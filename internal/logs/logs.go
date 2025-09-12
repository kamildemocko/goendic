package logs

import (
	"log/slog"
	"os"
	"path/filepath"
)

var logFile *os.File

func InitLogger(path string) error {
	logDir := filepath.Dir(path)
	_, notfound := os.Stat(logDir)
	if notfound != nil {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	var err error
	logFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
