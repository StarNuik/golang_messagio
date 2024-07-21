package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/message"
	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/starnuik/golang_messagio/internal/stream"
)

var messages *model.MessagesModel
var messageCreated *stream.DbMessageCreated

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

	err = messages.Insert(context.TODO(), msg)
	cmd.ServerErrorResponse(err, c)

	err = messageCreated.Publish(context.TODO(), msg.Id)
	cmd.ServerErrorResponse(err, c)

	c.JSON(http.StatusCreated, msg)
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func main() {
	db, err := internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ServerError(err)
	defer db.Close()

	messages = model.NewMessagesModel(db)

	messageCreated = stream.NewDbMessageCreated(os.Getenv("SERVICE_KAFKA_URL"), 10e3)
	defer messageCreated.Close()

	router := gin.Default()
	router.GET("/healthcheck", healthcheck)
	router.POST("/message", postMessageRequest)

	router.Run("0.0.0.0:8080")
}
