package message

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/starnuik/golang_messagio/internal/api"
	"github.com/starnuik/golang_messagio/internal/model"
)

func Validate(req api.MessageRequest) (model.Message, error) {
	msg := model.Message{}

	if len(req.Content) <= 0 {
		return msg, fmt.Errorf("zero-length content")
	}

	if len(req.Content) > 1024 {
		return msg, fmt.Errorf("content is too long, max size is 1024 bytes")
	}

	id, err := uuid.NewV4()
	if err != nil {
		return msg, fmt.Errorf("could not generate a uuid")
	}

	return model.Message{
		Id:      id,
		Created: time.Now().UTC(),
		Content: req.Content,
	}, nil
}

func Process(msg model.Message) model.Message {
	msg.IsProcessed = true
	return msg
}
