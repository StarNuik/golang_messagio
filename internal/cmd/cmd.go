package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal/api"
)

// panics on err != nil
func PanicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// sends an http.BadRequest
func ErrorResponse(err error, c *gin.Context) {
	if err != nil {
		errorResponse(err, c, http.StatusBadRequest)
		log.Println(err)
	}
}

func errorResponse(err error, c *gin.Context, status int) {
	res := api.ErrorResponse{
		Status:      status,
		Description: err.Error(),
	}
	c.JSON(res.Status, res)
}
