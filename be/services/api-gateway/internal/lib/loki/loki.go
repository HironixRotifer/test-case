package loki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	"net/http"
	"time"
)

var (
	urlToServerLoki string
	job             string
)

func init() {
	cfg := config.MustLoadConfig("./services/api-gateway/.env")

	urlToServerLoki = cfg.LokiURL
	job = cfg.Job
}

// LogEntry структура для отправки логов
type logEntry struct {
	Streams []stream `json:"streams"`
}

// Stream представляет собой поток логов
type stream struct {
	Labels  string  `json:"labels"`
	Entries []entry `json:"entries"`
}

// entry представляет собой саму запись лога
type entry struct {
	Timestamp time.Time `json:"ts"`   // Время записи
	Line      string    `json:"line"` // Сообщение
}

// generateLogEntry формирует лог
func (c *customLoki) generateLogEntry(message string) logEntry {

	entryLog := entry{
		Timestamp: time.Now(),
		Line:      message,
	}

	job := fmt.Sprintf("{job=\"%s\"}", job)

	Stream := stream{
		Labels:  job,
		Entries: []entry{entryLog},
	}

	logEntry := logEntry{
		Streams: []stream{Stream},
	}

	return logEntry
}

// sendLogToLoki отправляет лог на сервер Loki
func (c *customLoki) sendLogToLoki(logEntry logEntry) error {
	payload, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	req, err := http.NewRequest("POST", urlToServerLoki, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
