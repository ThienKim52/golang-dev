package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ThienKim52/golang-dev/docs"
	"github.com/ThienKim52/golang-dev/internal/handler"
	"github.com/ThienKim52/golang-dev/internal/repository"
	"github.com/ThienKim52/golang-dev/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var swaggerOnce sync.Once

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

func New(gin *gin.Engine, config *Config, redisClient *redis.Client) Engine {
	app := &engine{
		app:         gin,
		config:      config,
		redisClient: redisClient,
	}
	app.initRoutes()
	return app
}

// Start server
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.config.AppPort))
}

// ServeHTTP implements http.Handler interface
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initHandlers() *allHandlers{
	linkRepo := repository.NewRedisLinkRepository(e.redisClient)

	// Health check endpoint with Redis check
	healthCheckSvc := service.NewHealthCheck(linkRepo, e.config.ServiceName, e.config.InstanceID)
	healthCheckSvcHdl := handler.NewHealthCheck(healthCheckSvc)
	// Genpass endpoint
	genPassSvc := service.NewGenPass()
	genPassHdl := handler.NewGenPass(genPassSvc)
	// URL shortening endpoint
	linkSvc := service.NewLinkService(linkRepo, genPassSvc)
	linkHdl := handler.NewLinkHandler(linkSvc)

	return &allHandlers{
		healthCheckSvcHdl: healthCheckSvcHdl,
		genPassHdl:        genPassHdl,
		linkHdl:           linkHdl,
	}
}

type allHandlers struct {
	healthCheckSvcHdl handler.HealthCheckHandler
	genPassHdl        handler.GenPass
	linkHdl           handler.LinkHandler
}

// Initialize routes
func (e *engine) initRoutes() {
	// Initialize repository with Redis client
	allHandlers := e.initHandlers()
	swaggerOnce.Do(func() {
		docs.SwaggerInfo.BasePath = e.config.BasePath
	})
	e.app.GET("/health-check", allHandlers.healthCheckSvcHdl.GetResponse)

	e.app.GET("/genpass", allHandlers.genPassHdl.GeneratePassword)

	v1Routes := e.app.Group("/v1")
	{
		// Link-related
		v1Routes.POST("/links/shorten", allHandlers.linkHdl.ShortenURL)
		v1Routes.GET("/links/redirect/:code", allHandlers.linkHdl.Redirect)
	}

	// Swagger documentation route
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
