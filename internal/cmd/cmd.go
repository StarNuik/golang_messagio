package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal/api"
)

// calls os.Exit on err != nil
func ExitIf(err error) {
	if err != nil {
		// todo: os.Exit doesnt let the go runtime to run defer-red code
		log.Fatalln("ERROR: ", err)
	}
}

// calls os.Exit on err != nil, send an http.InternalServerError
func RespondAndExitIf(err error, c *gin.Context) {
	if err != nil {
		errorResponse(err, c, http.StatusInternalServerError)
		ExitIf(err)
	}
}

// sends an http.BadRequest, doesn't call os.Exit
func RespondIf(err error, c *gin.Context) {
	if err != nil {
		errorResponse(err, c, http.StatusBadRequest)
	}
}

func errorResponse(err error, c *gin.Context, status int) {
	res := api.ErrorResponse{
		Status:      status,
		Description: err.Error(),
	}
	c.JSON(res.Status, res)
}
