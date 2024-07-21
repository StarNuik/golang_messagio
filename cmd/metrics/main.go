package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal/model"
	// "github.com/starnuik/golang_messagio/pkg/lib"
)

var metrics *model.MetricsModel

func checkError(err error) {
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}
}

func getMetrics(c *gin.Context) {
	metrics, err := metrics.Get()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		checkError(err)
	}

	c.IndentedJSON(http.StatusOK, metrics)
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {
	var err error
	metrics, err = model.NewMetricsModel(os.Getenv("SERVICE_POSTGRES_URL"))
	checkError(err)

	router := gin.Default()

	router.GET("/healthcheck", healthcheck)
	router.GET("/metrics", getMetrics)

	router.Run("0.0.0.0:8080")
}
