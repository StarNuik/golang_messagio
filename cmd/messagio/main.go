package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
		if err != nil {
			log.Println(err)
			continue
		}

		go processMessage(msg)
	}
}

func main() {
	defer log.Println("cleaned up")

	var err error
	metrics, err = model.NewMetricsModel(context.Background(), postgresUrl)
	cmd.PanicIf(err)
	defer metrics.Close()

	messages, err = model.NewMessagesModel(context.Background(), postgresUrl)
	cmd.PanicIf(err)
	defer messages.Close()

	messageCreated, err = stream.NewDbMessageCreated(kafkaUrl, 10e3)
	cmd.PanicIf(err)
	defer messageCreated.Close()

	go messageReader()

	router := gin.Default()

	router.GET("/healthcheck", getHealthcheck)
	router.POST("/message", postMessage)
	router.GET("/query/metrics", getMetrics)
	router.GET("/query/message", getMessage)

	router.Run("0.0.0.0:8080")
}
