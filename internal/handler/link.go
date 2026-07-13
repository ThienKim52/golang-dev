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
type shortenURL struct {
	service service.LinkService
}

// NewLinkHandler creates a new link handler
func NewLinkHandler(service service.LinkService) *shortenURL {
	return &shortenURL{
		service: service,
	}
}
type req struct {
	URL string `json:"url" binding:"required"`
	Exp time.Duration  `json:"exp" binding:"required"`
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
func (h *shortenURL) ShortenURL(c *gin.Context) {
	req := &req{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	code, err := h.service.ShortenURL(c, req.URL, req.Exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shorten URL generated successfully!",
		"code":    code,
	})
}
