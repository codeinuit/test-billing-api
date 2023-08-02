package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlers struct {
}

func (h handlers) healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
