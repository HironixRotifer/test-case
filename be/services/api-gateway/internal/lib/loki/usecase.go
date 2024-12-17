package loki

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type customLoki struct {
	level string // уровень логирования
}

// Уровень info для логирования
func Info() *customLoki {
	return &customLoki{level: "info"}
}

// Уровень error для логирования
func Error() *customLoki {
	return &customLoki{level: "error"}
}

// Уровень debug для логирования
func Debug() *customLoki {
	return &customLoki{level: "debug"}
}

// Msgf логирует сообщение
func (c *customLoki) Msgf(format string, args ...interface{}) {
	const op = "Loki.Log"

	message := fmt.Sprintf(format, args...)

	logEntry := c.generateLogEntry(message)

	err := c.sendLogToLoki(logEntry)
	if err != nil {
		log.Error().Err(err).Msgf("%s failed to log ", op)
	}

}
