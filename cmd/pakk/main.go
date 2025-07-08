package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func fatal(logger *slog.Logger, message string, attrs ...any) {
	logger.Error(message, attrs...)
	os.Exit(1)
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	out := filepath.Join(root, ".pakk")
	logDir := filepath.Join(out, "logs")

	_ = os.MkdirAll(logDir, 0755)

	logFile, err := os.Create(filepath.Join(
		logDir,
		fmt.Sprintf("%s.log", time.Now().Format("20060102150405")),
	))
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewTextHandler(
		io.MultiWriter(os.Stdout, logFile),
		nil,
	))

	project, err := parseConfig(&parseConfigOptions{
		RootDir: ".",
		LogFile: logFile,
	})
	if err != nil {
		fatal(logger, "failed to parse config",
			slog.Any("error", err),
		)
	}
	defer project.Cleanup()

	if err := project.Build(); err != nil {
		fatal(logger, "failed to build project",
			slog.Any("error", err),
		)
	}
}
