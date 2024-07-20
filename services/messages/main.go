package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/lib"
)

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

	port := lib.EnvOr("MESSAGIO_PORT", "8080")
	router.Run("0.0.0.0:" + port)
}
