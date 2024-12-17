package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	shutDownTimeout = 10 * time.Second
)

type ServerHTTP struct {
	port   int
	server *http.Server

	Router *gin.Engine
}

// NewServerHTTP - создаёт новый http сервер
func NewServerHTTP(port int) *ServerHTTP {
	r := gin.Default()
	addr := fmt.Sprintf(":%v", port)
	server := &http.Server{Addr: addr, Handler: r}

	return &ServerHTTP{
		port:   port,
		server: server,
		Router: r,
	}
}

// Start запускает сервер
func (h *ServerHTTP) Start() {
	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("failed to start http server")
			}
		}
	}()

	log.Info().Msg("http server is started!")
}

// Stop gracefully stops the server
func (h *ServerHTTP) Stop(wg *sync.WaitGroup) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
		defer cancel()
		defer wg.Done()

		if err := h.server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server forced to shutdown")
		} else {
			log.Info().Msg("HTTP stopped")
		}
	}()
}
