package handler

import (
	"net/http"
	"time"

	"github.com/ThienKim52/golang-dev/internal/service"
	"github.com/gin-gonic/gin"
)

type LinkHandler interface {
	ShortenURL(c *gin.Context)
}

// LinkHandler handles link-related requests
type linkHandler struct {
	service service.LinkService
}

// NewLinkHandler creates a new link handler
func NewLinkHandler(service service.LinkService) *linkHandler {
	return &linkHandler{
		service: service,
	}
}

// ShortenURL handles POST /v1/links/shorten
// @Summary Shorten a URL
// @Description Creates a short code for a given URL
// @Tags links
// @Accept json
// @Produce json
// @Param url body string true "URL to shorten"
// @Param exp body int true "Expiration time in seconds"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *linkHandler) ShortenURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
		Exp int64  `json:"exp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.service.ShortenURL(c.Request.Context(), req.URL, time.Duration(req.Exp)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": "Shorten URL generated successfully!",
	})
}
