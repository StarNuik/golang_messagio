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
var messages *model.MessagesModel
var workloads *model.WorkloadsModel

func getMetrics(c *gin.Context) {
	metrics, err := metrics.Get(context.TODO())
	cmd.ServerErrorResponse(err, c)

	c.IndentedJSON(http.StatusOK, metrics)
}

func getMessage(c *gin.Context) {
	panic("todo")
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {
	db, err := internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ServerError(err)
	defer db.Close()

	metrics = model.NewMetricsModel(db)
	messages = model.NewMessagesModel(db)
	workloads = model.NewWorkloadsModel(db)

	router := gin.Default()

	router.GET("/healthcheck", healthcheck)
	router.GET("/query/metrics", getMetrics)
	router.GET("/query/message", getMessage)

	router.Run("0.0.0.0:8080")
}
