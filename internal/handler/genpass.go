package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ThienKim52/golang-dev/internal/service"
)

const passwordLength = 12

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPass struct {
	svc service.GenPass
}

func NewGenPass(svc service.GenPass) GenPass {
	return &genPass{svc: svc}
}

func (g *genPass) GeneratePassword(c *gin.Context) {
	// call service
	pass, err := g.svc.GeneratePassword(passwordLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// return response
	c.JSON(http.StatusOK, gin.H{"password": pass})
}
