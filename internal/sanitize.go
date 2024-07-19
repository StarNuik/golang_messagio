package internal

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func MessageFromReq(req MessageRequest) (Message, error) {
	if len(req.Content) == 0 {
		return Message{}, fmt.Errorf("ERROR: empty message content")
	}
	len := min(len(req.Content), 4096)
	content := req.Content[:len]

	id, err := uuid.NewV4()
	if err != nil {
		return Message{}, err
	}

	return Message{
		Id:        id,
		Created:   time.Now().UTC(),
		Content:   content,
		Processed: false,
	}, nil
}
