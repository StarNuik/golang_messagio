package api

import "github.com/gofrs/uuid/v5"

type ErrorResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

// '/message' GET
type MessageRequest struct {
	Content string `json:"content"`
}

// '/query/message' GET
type MessageQueryRequest struct {
	Id uuid.UUID `json:"id"`
}
