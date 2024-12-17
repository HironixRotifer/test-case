package main

import (
	"fmt"
	pp "gateway/cmd/app/prometheus"
	authGRPClient "gateway/delivery/grpc/authGRPClient"
	serve "gateway/delivery/http"
	"gateway/delivery/http/routes"
	"gateway/internal/config"
	"os"
	"os/signal"
	"syscall"

	"sync"

	"github.com/rs/zerolog/log"
)

var wg sync.WaitGroup

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	cfg := config.MustLoadConfig("./services/api-gateway/.env")
	fmt.Println(cfg)

	server := serve.NewServerHTTP(8080)

	adminGroup := server.Router.Group("/api-v1/")
	userGroup := server.Router.Group("/")

	authClient := authGRPClient.NewAuthGRPClient("")
	routes.UserGroup(userGroup, authClient)
	routes.AdminGroup(adminGroup)

	pp.Start()
	// Using middleware
	// beforeAuthorization.Use(sessions.Sessions("session-name", store))
	// afterAuthorization.Use(sessions.Sessions("session-name", store))
	// afterAuthorization.Use(middlewares.AuthSession)

	server.Start()

	<-signalChan
	log.Info().Msg("server graceful stop started")

	wg.Add(1)
	server.Stop(&wg)
	wg.Wait()
}
