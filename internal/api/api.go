package api

import (
	"net/http"

	_ "github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/handler"
	"github.com/ThienKim52/golang-dev/internal/repository"
	"github.com/ThienKim52/golang-dev/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app         *gin.Engine
	config      *Config
	redisClient *redis.Client
}

func NewEngine(config *Config, redisClient *redis.Client) Engine {
	app := &engine{
		app:         gin.Default(),
		config:      config,
		redisClient: redisClient,
	}
	app.initRoutes()
	return app
}

func (e *engine) Start() error {
	return e.app.Run("localhost:8080")
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initRoutes() {
	// Initialize repository with Redis client
	linkRepo := repository.NewRedisLinkRepository(e.redisClient)

	// Health check endpoint with Redis check
	healthCheckSvc := service.NewHealthCheck(linkRepo)
	healthCheckSvcHdl := handler.NewHealthCheck(healthCheckSvc)

	// Middleware to inject service name and instance ID into context
	middleware := func(c *gin.Context) {
		c.Set("service_name", e.config.ServiceName)
		c.Set("instance_id", e.config.InstanceID)
		c.Next()
	}

	// URL shortening endpoint
	genPassSvc := service.NewGenPass()
	linkSvc := service.NewLinkService(linkRepo, genPassSvc)
	linkHdl := handler.NewLinkHandler(linkSvc)
	e.app.POST("/v1/links/shorten", linkHdl.ShortenURL)
	e.app.GET("/health-check", middleware, healthCheckSvcHdl.GetResponse)

	// Swagger documentation route
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
