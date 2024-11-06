package http

import (
	"errors"
	"fmt"
	"gateway/http/routes"

	"math/rand/v2"

	"github.com/gin-gonic/gin"
)

const (
	defaultHTTPPort = 8080
)

type Server struct {
	addr string
	port int

	g *gin.Engine
}

// New creates a new gin Server.
// have ip address and port.
func NewServer(addr string, opts ...Option) *Server {
	var options options

	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			// log.Err(err).Msg("error starting server")
			return nil
		}
	}

	var port int
	if options.port == nil {
		port = defaultHTTPPort
	} else {
		if *options.port == 0 {
			port = randomPort()
		} else {
			port = *options.port
		}
	}

	router := gin.Default()

	// Initializing group router
	adminGroup := router.Group("/api-v1/")
	userGroup := router.Group("/")
	// Using middleware
	// beforeAuthorization.Use(sessions.Sessions("session-name", store))
	// afterAuthorization.Use(sessions.Sessions("session-name", store))
	// afterAuthorization.Use(middlewares.AuthSession)

	// Routing
	routes.UserGroup(userGroup)
	routes.AdminGroup(adminGroup)

	return &Server{
		addr: addr,
		port: port,

		g: router,
	}
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%v", s.addr, s.port)

	err := s.g.Run(addr)
	if err != nil {
		// log
		return err
	}

	// log

	return nil
}

// TODO
// Stop gracefully stops the server
func (s *Server) Stop() error {

	return nil
}

type options struct {
	port *int
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		options.port = &port
		return nil
	}
}

func randomPort() int {
	return rand.IntN(9000)
}
