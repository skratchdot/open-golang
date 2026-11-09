package opengolang_test

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("cannot get current working dir: %w", err))
	}
	logsDir := filepath.Join(wd, "logs")
	_ = os.RemoveAll(logsDir)
	slog.Info("reseting logs directory", slog.String("path", logsDir))
	m.Run()
}
