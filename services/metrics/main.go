package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/lib"
)

func getMetrics(c *gin.Context) {
	metrics := lib.Metrics{}

	c.IndentedJSON(http.StatusOK, metrics)
}

func main() {
	router := gin.New()

	router.GET("/metrics", getMetrics)

	port := lib.EnvOr("MESSAGIO_PORT", "8080")
	router.Run("0.0.0.0:" + port)
}
