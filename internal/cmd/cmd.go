package cmd

import (
	"log"

	"github.com/gin-gonic/gin"
)

// panics on err != nil
func PanicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func ErrorResponse(c *gin.Context, err error, desc string, status int) {
	log.Println("ERROR:", err)
	c.String(status, "%s", desc)
}
