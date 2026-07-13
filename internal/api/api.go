package api

import (
	"fmt"
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

// interface Engine
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app         *gin.Engine
	config      *Config
	redisClient *redis.Client
}

// Constructor
func NewEngine(config *Config, redisClient *redis.Client) Engine {
	app := &engine{
		app:         gin.Default(),
		config:      config,
		redisClient: redisClient,
	}
	app.initRoutes()
	return app
}

// Start server
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.config.Port))
}

// ServeHTTP implements http.Handler interface
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// Initialize routes
func (e *engine) initRoutes() {
	// Initialize repository with Redis client
	linkRepo := repository.NewRedisLinkRepository(e.redisClient)

	// Health check endpoint with Redis check
	healthCheckSvc := service.NewHealthCheck(linkRepo, e.config.ServiceName, e.config.InstanceID)
	healthCheckSvcHdl := handler.NewHealthCheck(healthCheckSvc)

	// URL shortening endpoint
	genPassSvc := service.NewGenPass()
	linkSvc := service.NewLinkService(linkRepo, genPassSvc)
	linkHdl := handler.NewLinkHandler(linkSvc)
	e.app.POST("/v1/links/shorten", linkHdl.ShortenURL)
	e.app.GET("/health-check", healthCheckSvcHdl.GetResponse)

	// Swagger documentation route
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
