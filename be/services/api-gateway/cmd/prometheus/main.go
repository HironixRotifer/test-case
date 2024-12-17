package pp

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	metric "gateway/internal/lib/prometheus"

	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	shutDownTimeout = 10 * time.Second
)

var wg sync.WaitGroup

func Start() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	registry := prometheus.NewRegistry()
	memory := metric.NewMetrics(registry)

	ctx, cancel := context.WithCancel(context.Background())

	go memory.CollectMemUsageMetric(ctx, &wg)
	server := startNewPrometheuServer(registry)

	<-signalChan
	log.Info().Msg("prometheus graceful stop started")

	wg.Add(1)
	go stopPrometheuServer(server)
	cancel()

	wg.Wait()
}

func startNewPrometheuServer(registry *prometheus.Registry) *http.Server {
	server := &http.Server{Addr: ":9092"}
	http.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{Registry: registry},
	))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("failed to start http server")
			}
		}
	}()

	log.Info().Msg("prometheus server is started!")

	return server
}

func stopPrometheuServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)

	defer cancel()
	defer wg.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("HTTP stopped")
	}

}
