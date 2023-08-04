package main

import (
	"codeinuit/test-billing-api/pkg/log/logrus"

	"github.com/gin-gonic/gin"
)

type BillingAPI struct {
}

func main() {
	r := gin.Default()
	l := logrus.NewLogrusLogger()
	h := &handlers{log: l}

	r.GET("/health", h.healthcheck)

	r.Run()
}
