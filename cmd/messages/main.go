package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/segmentio/kafka-go"
	"github.com/starnuik/golang_messagio/internal/model"
)

type MessageReq struct {
	Content string `json:"content"`
}

var messages *model.MessagesModel

func checkError(err error) {
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}
}

func toMessage(req MessageReq) (model.Message, error) {
	msg := model.Message{}

	if len(req.Content) <= 0 {
		return msg, fmt.Errorf("zero-length content")
	}

	len := min(len(req.Content), 4096)
	id, err := uuid.NewV4()
	if err != nil {
		return msg, fmt.Errorf("could not generate a uuid")
	}

	return model.Message{
		Id:      id,
		Created: time.Now().UTC(),
		Content: req.Content[:len],
	}, nil
}

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
	var req MessageReq

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	msg, err := toMessage(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = messages.Insert(msg)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		checkError(err)
	}

	err = publishMessageCreated(msg.Id)
	if err != nil {
		checkError(err)
	}

	c.JSON(http.StatusCreated, msg)
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {
	var err error
	messages, err = model.NewMessagesModel(os.Getenv("SERVICE_POSTGRES_URL"))
	checkError(err)

	router := gin.Default()
	router.GET("/healthcheck", healthcheck)
	router.POST("/message", postMessageRequest)

	router.Run("0.0.0.0:8080")
}
