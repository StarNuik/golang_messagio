package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/starnuik/golang_messagio/internal/api"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/message"
	"github.com/starnuik/golang_messagio/internal/model"
)

func getHealthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func postMessage(c *gin.Context) {
	var req api.MessageRequest

	err := c.BindJSON(&req)
	if err != nil {
		cmd.ErrorResponse(err, c)
		return
	}

	msg, err := message.Validate(req)
	if err != nil {
		cmd.ErrorResponse(err, c)
		return
	}

	err = messages.Insert(context.TODO(), msg)
	cmd.PanicAndRespond(err, c)

	err = messageCreated.Publish(context.TODO(), msg)
	cmd.PanicAndRespond(err, c)

	c.JSON(http.StatusCreated, msg)
}

func processMessage(msg model.Message) {
	msg = message.Process(msg)

	err := messages.Update(context.TODO(), msg)
	cmd.Panic(err)
}

func getMetrics(c *gin.Context) {
	metrics, err := metrics.Get(context.TODO())
	cmd.PanicAndRespond(err, c)

	c.IndentedJSON(http.StatusOK, metrics)
}

func getMessage(c *gin.Context) {
	var req api.MessageQueryRequest

	err := c.BindJSON(&req)
	if err != nil {
		cmd.ErrorResponse(err, c)
		return
	}

	msg, err := messages.Get(context.TODO(), req.Id)
	if errors.Is(err, pgx.ErrNoRows) {
		c.Status(http.StatusNotFound)
		return
	}
	cmd.PanicAndRespond(err, c)

	c.IndentedJSON(http.StatusOK, msg)
}
