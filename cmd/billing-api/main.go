package main

import "github.com/gin-gonic/gin"

type BillingAPI struct {
}

func main() {
	r := gin.Default()
	h := &handlers{}

	r.GET("/health", h.healthcheck)
	r.Run()
}
