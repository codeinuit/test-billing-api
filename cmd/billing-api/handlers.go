package main

import (
	"net/http"

	logger "github.com/codeinuit/test-billing-api/pkg/log"

	"github.com/gin-gonic/gin"
)

type handlers struct {
	log logger.Logger
}

// healthcheck works as a ping and returns a OK status
func (h handlers) healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
