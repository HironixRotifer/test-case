package handlers

import (
	"context"
	"net/http"
	ssov1 "protos/auth"
	"time"

	authGRPClient "gateway/delivery/grpc/authGRPClient"
	"gateway/internal/lib/loki"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	grpcClient authGRPClient.AuthClient
	validator  *validator.Validate
}

func NewAuthorizationHandler(client authGRPClient.AuthClient) *Handler {
	return &Handler{
		grpcClient: client,
		validator:  validator.New(),
	}
}

// LoginUser - вызывает gRPC обработчик авторизации пользователя.
// Успешная авторизация возвращает Json Web Token пользователя.
// В случае провала возвращает и логирует ошибку.
func (h *Handler) HandleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		var logger = log.With().Str("service", "authorization").Str("method", "Login").Logger()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		validationErr := h.validator.Struct(req)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.grpcClient.Login(ctx, &ssov1.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})

		if err != nil {
			logger.Error().Err(err).Msgf("error to call method of login user %v", resp.Error)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"JwToken": resp.Token})
	}
}

// HandleRegisterNewUser - обработчик вызова gRPC регистрации пользователя.
// Успешное выполнение возвращает ID зарегистрированного пользователя.
func (h *Handler) HandleRegisterNewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		var logger = log.With().Str("service", "authorization").Str("method", "Register").Logger()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		validationErr := h.validator.Struct(req)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.grpcClient.Register(ctx, &ssov1.RegisterRequest{
			Email:       req.Email,
			Password:    req.Password,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			PhoneNumber: req.PhoneNumber,
		})

		if err != nil {
			logger.Error().Err(err).Msgf("error to call method to register user %v", resp.Error)
		}

		c.JSON(http.StatusCreated, gin.H{"UserID": resp.UserId})
	}
}

// HandleIsAdmin вызывает gRPC обработчик проверки статуса администратора у пользователя.
// Успешная проверка возвращает.
// В случае провала возвращает событие False и логирует ошибку.
func (h *Handler) HandleIsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req isAdminRequest

		var logger = log.With().
			Str("service", "authorization").
			Str("method", "IsAdmin").
			Logger()

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}

		validationErr := h.validator.Struct(req)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bad request",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.grpcClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
			UserId: req.ID,
		})

		if err != nil {
			logger.Error().Err(err).Msgf("error to call method to require admin rights %v", resp.Error)
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"isAdmin": resp.IsAdmin,
		})
	}
}

func (h *Handler) Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Hello": "Bro!"})
		// loki.SendLogToLoki("TEST")
		// slice := make([]byte, 1024)
		// fmt.Println(slice)
		// l := loki.SetupLogger(nil)

		loki.Info().Msgf("SOSAT` AMERICA_3")

		// prometheus.NewMetrics()
	}
}
