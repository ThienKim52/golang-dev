package handler

import (
	"net/http"

	"github.com/ThienKim52/golang-dev/internal/service"
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler handles health check requests
type HealthCheckHandler struct {
	service service.HealthCheckService
	serviceName string
	instanceID string
}

// NewHealthCheck creates a new health check handler
func NewHealthCheck(service service.HealthCheckService, serviceName, instanceID string) *HealthCheckHandler {
	return &HealthCheckHandler{
		service:     service,
		serviceName: serviceName,
		instanceID:  instanceID,
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

	message, err := h.service.GetHealthCheck(c.Request.Context(), h.serviceName, h.instanceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis connection failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"service_name": h.serviceName,
		"instance_id":  h.instanceID,
	})
}
