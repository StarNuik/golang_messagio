package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	status      int
	description string
}

// calls os.Exit on err != nil
func ServerError(err error) {
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}
}

// calls os.Exit on err != nil
func ServerErrorResponse(err error, c *gin.Context) {
	if err != nil {
		res := ErrorResponse{
			status:      http.StatusInternalServerError,
			description: err.Error(),
		}
		c.JSON(res.status, res)
		ServerError(err)
	}
}
