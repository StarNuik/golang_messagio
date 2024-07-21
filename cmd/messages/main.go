package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/segmentio/kafka-go"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/message"
	"github.com/starnuik/golang_messagio/internal/model"
)

var messages *model.MessagesModel

func publishMessageCreated(msgId uuid.UUID) error {
	cfg := kafka.WriterConfig{
		Brokers: []string{os.Getenv("SERVICE_KAFKA_URL")},
		Topic:   "db.message.created",
	}
	w := kafka.NewWriter(cfg)
	//todo: fails on the first message publish
	w.AllowAutoTopicCreation = true

	kMsg := kafka.Message{Key: nil, Value: msgId.Bytes()}
	err := w.WriteMessages(context.TODO(), kMsg)
	return err
}

func postMessageRequest(c *gin.Context) {
	var req internal.MessageRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	msg, err := message.Validate(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = messages.Insert(msg)
	cmd.ServerErrorResponse(err, c)

	err = publishMessageCreated(msg.Id)
	cmd.ServerErrorResponse(err, c)

	c.JSON(http.StatusCreated, msg)
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {
	var err error
	messages, err = model.NewMessagesModel(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ServerError(err)

	router := gin.Default()
	router.GET("/healthcheck", healthcheck)
	router.POST("/message", postMessageRequest)

	router.Run("0.0.0.0:8080")
}
