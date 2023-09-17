package loggy

import (
	"fmt"
	"log/slog"
	"strings"
)

func checkLevel(log string, level slog.Level) bool {
	return strings.Contains(log, fmt.Sprintf("level=%s", level.String())) ||
		strings.Contains(log, fmt.Sprintf("\"level\":\"%s\"", level.String()))
}
