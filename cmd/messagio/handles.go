package main

import (
	"context"
	"errors"
	"log"
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
		cmd.ErrorResponse(c, err, "could not parse json body", http.StatusBadRequest)
		return
	}

	msg, err := message.Validate(req)
	if err != nil {
		cmd.ErrorResponse(c, err, "invalid request data", http.StatusBadRequest)
		return
	}

	err = messages.Insert(context.Background(), msg)
	if err != nil {
		cmd.ErrorResponse(c, err, "database insert error", http.StatusInternalServerError)
		return
	}

	err = messageCreated.Publish(context.Background(), msg)
	if err != nil {
		cmd.ErrorResponse(c, err, "broker message publishing error", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func processMessage(msg model.Message) {
	msg = message.Process(msg)

	err := messages.UpdateIsProcessed(context.Background(), msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func getMetrics(c *gin.Context) {
	metrics, err := metrics.Get(context.Background())
	if err != nil {
		cmd.ErrorResponse(c, err, "database read error", http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, metrics)
}

func getMessage(c *gin.Context) {
	var req api.MessageQueryRequest

	err := c.BindJSON(&req)
	if err != nil {
		cmd.ErrorResponse(c, err, "could not parse json body", http.StatusBadRequest)
		return
	}

	msg, err := messages.Get(context.Background(), req.Id)
	if errors.Is(err, pgx.ErrNoRows) {
		cmd.ErrorResponse(c, err, "could not find such entry", http.StatusNotFound)
		return
	}
	if err != nil {
		cmd.ErrorResponse(c, err, "database read error", http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, msg)
}
