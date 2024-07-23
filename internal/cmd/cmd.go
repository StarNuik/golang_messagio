package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/starnuik/golang_messagio/internal/api"
)

// calls os.Exit on err != nil
func Panic(err error) {
	if err != nil {
		// todo: os.Exit doesnt let the go runtime to run defer-red code
		log.Fatalln("ERROR: ", err)
	}
}

// calls os.Exit on err != nil, send an http.InternalServerError
func PanicAndRespond(err error, c *gin.Context) {
	if err != nil {
		errorResponse(err, c, http.StatusInternalServerError)
		Panic(err)
	}
}

// sends an http.BadRequest, doesn't call os.Exit
func ErrorResponse(err error, c *gin.Context) {
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
