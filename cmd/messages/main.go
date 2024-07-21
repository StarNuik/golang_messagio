package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/starnuik/golang_messagio/internal"
)

type MessageReq struct {
	Content string `json:"content"`
}

type Message struct {
	Id      uuid.UUID
	Created time.Time
	Content string
}

var sql *pgxpool.Pool

func checkError(err error) {
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}
}

func healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func postMessageRequest(c *gin.Context) {
	var req MessageReq

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	msg, err := newMessage(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}

	err = insertMessage(msg)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Fatalf("ERROR: %v", err)
	}

	c.JSON(http.StatusCreated, msg)
}

func main() {
	var err error
	sql, err = internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	checkError(err)

	router := gin.Default()
	router.GET("/healthcheck", healthcheck)
	router.POST("/message", postMessageRequest)

	router.Run("0.0.0.0:8080")
}
