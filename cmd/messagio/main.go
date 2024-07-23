package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/starnuik/golang_messagio/internal/stream"
)

var (
	metrics        *model.MetricsModel
	messages       *model.MessagesModel
	messageCreated *stream.DbMessageCreated
	postgresUrl    = os.Getenv("SERVICE_POSTGRES_URL")
	kafkaUrl       = os.Getenv("SERVICE_KAFKA_URL")
)

func messageReader() {
	for {
		msg, err := messageCreated.Read(context.TODO())
		cmd.Panic(err)
		go processMessage(msg)
	}
}

func main() {
	db, err := internal.NewSqlPool(postgresUrl)
	cmd.Panic(err)
	defer db.Close()

	metrics = model.NewMetricsModel(db)
	messages = model.NewMessagesModel(db)

	messageCreated = stream.NewDbMessageCreated(kafkaUrl, 10e3)
	defer messageCreated.Close()

	go messageReader()

	router := gin.Default()

	router.GET("/healthcheck", getHealthcheck)
	router.POST("/message", postMessage)
	router.GET("/query/metrics", getMetrics)
	router.GET("/query/message", getMessage)

	router.Run("0.0.0.0:8080")
}
