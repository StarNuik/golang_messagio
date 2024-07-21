package api

import (
	"github.com/gofrs/uuid/v5"
	"github.com/starnuik/golang_messagio/internal/model"
)

type ErrorResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// '/message' POST
type MessageRequest struct {
	Content string `json:"content"`
}

// '/query/message' GET
type MessageQueryRequest struct {
	Id uuid.UUID `json:"id"`
}

// '/query/message' GET
type MessageQueryResponse struct {
	Message   *model.Message   `json:"message"`
	Processed *model.Processed `json:"processed"`
}
