package loki

import (
	"bytes"
	"clientCorner/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loki struct {
	c *customLoki
}

type customLoki struct {
	lokiURL string // url to loki
	job     string // name of service
	level   string // logging level
}

func SetupLogger(cfg *config.Config) *loki {
	return &loki{
		c: &customLoki{
			job:     cfg.Job,
			lokiURL: cfg.LokiURL,
		},
	}
}

func (l *loki) Info() *customLoki {
	return &customLoki{level: "info"}
}

func (l *loki) Error() *customLoki {
	return &customLoki{level: "error"}
}

func (l *loki) Debug() *customLoki {
	return &customLoki{level: "debug"}
}

// LogEntry структура для отправки логов
type logEntry struct {
	Streams []stream `json:"streams"`
}

// Stream представляет собой поток логов
type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// Msgf форматирует сообщение и записывает его в лог
func (c *customLoki) Msgf(format string, args ...interface{}) {
	const op = "Loki.Log"

	msg := fmt.Sprintf(format, args...)

	entry := logEntry{
		Streams: []stream{
			{
				Stream: map[string]string{"job": c.job, "level": c.level},
				Values: [][]string{
					{timestamp(), msg},
				},
			},
		},
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Err(err).Msgf("%s Failed to marshal log entry", op)
		return
	}

	resp, err := http.Post(c.lokiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Msgf("%s Failed to send log to Loki: ", op)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("%s Failed log to Loki: received status %s", op, resp.Status)
	}
}

func timestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}
