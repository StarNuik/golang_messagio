package stream

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/starnuik/golang_messagio/internal/model"
)

type DbMessageCreated struct {
	*kafkaStream
}

func NewDbMessageCreated(brokerUrl string, messageSize int) (*DbMessageCreated, error) {
	const topic = "db.message.created"

	stream := &DbMessageCreated{
		newStream(brokerUrl, topic, messageSize),
	}

	err := stream.createTopic()
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (s *DbMessageCreated) Publish(ctx context.Context, msg model.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w := s.writer()
	km := kafka.Message{Key: msg.Id.Bytes(), Value: payload}
	err = w.WriteMessages(ctx, km)
	return err
}

func (s *DbMessageCreated) Read(ctx context.Context) (model.Message, error) {
	msg := model.Message{}

	r := s.reader()
	km, err := r.ReadMessage(ctx)
	if err != nil {
		return msg, err
	}

	err = json.Unmarshal(km.Value, &msg)
	return msg, err
}

func (s *DbMessageCreated) Close() error {
	return s.close()
}
