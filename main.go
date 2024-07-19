package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Content string `json:"content"`
}

type Metrics struct {
	MessagesReceived int `json:"messagesReceived"`
}

var metricsStore Metrics

func postMessage(c *gin.Context) {
	var msg Message

	err := c.BindJSON(&msg)
	if err != nil {
		return
	}

	metricsStore.MessagesReceived += 1

	fmt.Println("Received a message:", msg)
	c.JSON(http.StatusCreated, msg)
}

func getMetrics(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, metricsStore)
}

func main() {
	router := gin.Default()
	router.POST("/message", postMessage)
	router.GET("/metrics", getMetrics)

	router.Run("0.0.0.0:8080")
}
