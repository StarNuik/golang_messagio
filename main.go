package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang_messagio/internal"
	_ "github.com/joho/godotenv/autoload"
)

var (
	metricsStore internal.Metrics
	kafka        *internal.KafkaDriver
	sql          *internal.SqlDriver
)

func postMessage(c *gin.Context) {
	var msg internal.Message

	err := c.BindJSON(&msg)
	if err != nil {
		return
	}

	err = kafka.Emit("post", msg)
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

func setup() (cleanup func()) {
	var err error
	sql, err = internal.NewSqlDriver(os.Getenv("MESSAGIO_POSTGRES_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	kafka, err = internal.NewKafkaDriver(os.Getenv("MESSAGIO_KAFKA_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	cleanup = func() {
		sql.Cleanup()
		kafka.Cleanup()
	}
	return
}

func main() {
	cleanup := setup()
	defer cleanup()

	router := gin.Default()
	router.POST("/message", postMessage)
	router.GET("/metrics", getMetrics)

	router.Run("0.0.0.0:8080")
}

// func workloadProcessor() {}
