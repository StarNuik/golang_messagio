package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/starnuik/golang_messagio/lib"
)

type MessageReq struct {
	Content string `json:content`
}

type MessageValid struct {
	Id      uuid.UUID
	Created time.Time
	Content string
}

func newMessage(req MessageReq) (MessageValid, error) {
	return MessageValid{}, fmt.Errorf("not implemented")
}

func postMessage(c *gin.Context) {
	var req lib.MessageReq

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	log.Println("req received:", req)

	c.JSON(http.StatusCreated, req)
}

func main() {
	router := gin.Default()
	router.POST("/message", postMessage)

	router.Run("0.0.0.0:8080")
}
