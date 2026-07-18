package handler

import (
	"net/http"

	"github.com/ThienKim52/golang-dev/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct {
	service service.HealthCheckService
}

// NewHealthCheck creates a new health check handler
func NewHealthCheck(service service.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{
		service: service,
	}
}

// GetResponse handles GET /health-check
// @Summary Health check endpoint
// @Description Returns the health status of the service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /health-check [get]
func (h *HealthCheckHandler) GetResponse(c *gin.Context) {

	result, err := h.service.GetHealthCheck(c)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.HealthCheckHandler.GetResponse").Msg("Failed to get response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, result)
}
