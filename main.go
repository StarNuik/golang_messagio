package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang_messagio/internal"
	_ "github.com/joho/godotenv/autoload"
)

var (
	kafka *internal.KafkaDriver
	sql   *internal.SqlDriver
	// isDebug bool
)

func checkError(err error) {
	if err == nil {
		return
	}

	log.Fatalf("ERROR: %v", err)
}

func postMessage(c *gin.Context) {
	var req internal.MessageRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	msg, err := internal.MessageFromReq(req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = sql.InsertMessage(msg)
	checkError(err)

	err = kafka.EmitId("new", msg.Id)
	checkError(err)

	c.JSON(http.StatusCreated, msg)
}

func getMetrics(c *gin.Context) {
	metrics, err := sql.QueryMetrics()
	checkError(err)

	c.IndentedJSON(http.StatusOK, metrics)
}

func setup() (cleanup func()) {
	// isDebug = os.Getenv("MESSAGIO_DEBUG") == "1"

	var err error
	sql, err = internal.NewSqlDriver(os.Getenv("MESSAGIO_POSTGRES_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	kafka, err = internal.NewKafkaDriver("messages", os.Getenv("MESSAGIO_KAFKA_URL"))
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
