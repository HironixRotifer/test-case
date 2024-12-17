package routes

import (
	authGRPClient "gateway/delivery/grpc/authGRPClient"
	handlers "gateway/delivery/http/handlers/authorization"

	"github.com/gin-gonic/gin"
)

func UserGroup(rg *gin.RouterGroup, client authGRPClient.AuthClient) {

	auth := handlers.NewAuthorizationHandler(client)
	rg.GET("/test", auth.Test())
	rg.POST("/login", auth.HandleLogin())
	rg.POST("/signUp", auth.HandleRegisterNewUser())
	rg.POST("/isAdmin", auth.HandleIsAdmin())

}

func AdminGroup(rg *gin.RouterGroup) {

}
