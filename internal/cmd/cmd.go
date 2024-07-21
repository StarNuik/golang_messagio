package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal/api"
)

// calls os.Exit on err != nil
func ServerError(err error) {
	if err != nil {
		log.Fatalln("ERROR: ", err)
	}
}

// calls os.Exit on err != nil
func ServerErrorResponse(err error, c *gin.Context) {
	if err != nil {
		res := api.ErrorResponse{
			Status:      http.StatusInternalServerError,
			Description: err.Error(),
		}
		c.JSON(res.Status, res)
		ServerError(err)
	}
}
