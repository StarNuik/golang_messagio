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

func main() {
	var err error
	sql, err = internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	checkError(err)

	router := gin.Default()
	router.GET("/healthcheck", healthcheck)
	router.POST("/message", postMessageRequest)

	router.Run("0.0.0.0:8080")
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

func newMessage(req MessageReq) (Message, error) {
	if len(req.Content) <= 0 {
		return Message{}, fmt.Errorf("zero-length content")
	}

	len := min(len(req.Content), 4096)
	id, err := uuid.NewV4()
	if err != nil {
		return Message{}, fmt.Errorf("could not generate a uuid")
	}

	return Message{
		Id:      id,
		Created: time.Now().UTC(),
		Content: req.Content[:len],
	}, nil
}

func insertMessage(msg Message) error {
	tag, err := sql.Exec(
		context.Background(),
		"INSERT INTO messages (msg_id, msg_created, msg_content) VALUES ($1, $2, $3)",
		msg.Id, msg.Created, msg.Content)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		//todo: should this be a "stop the service" error or a plain log.error
		return fmt.Errorf("rowsAffected != 1")
	}
	return nil
}
