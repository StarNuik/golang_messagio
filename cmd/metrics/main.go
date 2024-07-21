package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/model"
)

var metrics *model.MetricsModel

func getMetrics(c *gin.Context) {
	metrics, err := metrics.Get(context.TODO())
	cmd.ServerErrorResponse(err, c)

	c.IndentedJSON(http.StatusOK, metrics)
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {
	db, err := internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ServerError(err)
	defer db.Close()

	metrics = model.NewMetricsModel(db)

	router := gin.Default()

	router.GET("/healthcheck", healthcheck)
	router.GET("/metrics", getMetrics)

	router.Run("0.0.0.0:8080")
}
