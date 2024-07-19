package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/lovoo/goka"
)

// todo: is this the correct way? it looks v hacky
type JsonCodec[T any] struct{}

func (_ *JsonCodec[T]) Encode(value interface{}) ([]byte, error) {
	data, err := json.Marshal(value)
	return data, err
}

func (_ *JsonCodec[T]) Decode(data []byte) (interface{}, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return t, err
}

type Message struct {
	Content string `json:"content"`
}

type Metrics struct {
	MessagesReceived int `json:"messagesReceived"`
}

var (
	metricsStore Metrics
	brokers                  = []string{"kafka:9092"}
	topic        goka.Stream = "messages"
)

func postMessage(c *gin.Context) {
	var msg Message

	err := c.BindJSON(&msg)
	if err != nil {
		return
	}

	emitter, err := goka.NewEmitter(brokers, topic, new(JsonCodec[Message]))
	if err != nil {
		log.Fatalln("ERROR: could not create a kafka/goka emitter:", err)
	}
	defer emitter.Finish()

	err = emitter.EmitSync("post", msg)
	if err != nil {
		log.Fatalln("ERROR: could not emit a kafka/goka message: ", err)
	}

	metricsStore.MessagesReceived += 1

	fmt.Println("Received a message:", msg)
	c.JSON(http.StatusCreated, msg)
}

func getMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, metricsStore)
}

func main() {
	cfg := goka.DefaultConfig()
	cfg.Version = sarama.V3_5_0_0
	goka.ReplaceGlobalConfig(cfg)

	router := gin.Default()
	router.POST("/message", postMessage)
	router.GET("/metrics", getMetrics)

	router.Run("0.0.0.0:8080")
}

// func workloadProcessor() {}
