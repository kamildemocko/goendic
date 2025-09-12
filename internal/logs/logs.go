package logs

import (
	"log/slog"
	"os"
	"path/filepath"
)

var logFile *os.File

func InitLogger() error {
	var err error

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(configDir, "goendic")
	_, notfound := os.Stat(logDir)
	if notfound != nil {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	logPath := filepath.Join(logDir, "words.log")
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
