package main

import (
	metric "clientCorner/internal/lib/prometheus"
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var wg *sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	reg := prometheus.NewRegistry()
	m := metric.NewMetrics(reg)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	server := &http.Server{Addr: ":8090"}

	wg.Add(1)
	go m.CollectMemUsageMetric(ctx, wg)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	cancel()
	server.Shutdown(context.Background())

	wg.Wait()
}
