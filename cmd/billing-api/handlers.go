package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlers struct {
}

func (h handlers) healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
