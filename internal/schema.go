package internal

import (
	"time"

	"github.com/gofrs/uuid"
)

type MessageCodec = JsonCodec[Message]

type MessageRequest struct {
	Content string `json:"content"`
}

type Message struct {
	Id        uuid.UUID
	Created   time.Time
	Content   string
	Processed bool
}

type Metrics struct {
	MessagesTotal     int `json:"messagesTotal"`
	MessagesLastDay   int `json:"messagesLastDay"`
	MessagesLastHour  int `json:"messagesLastHour"`
	ProcessedTotal    int `json:"processedTotal"`
	ProcessedLastDay  int `json:"processedLastDay"`
	ProcessedLastHout int `json:"processedLastHout"`
	// MessagesReceived int `json:"messagesReceived"`
}
