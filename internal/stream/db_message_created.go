package stream

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/segmentio/kafka-go"
)

type DbMessageCreated struct {
	broker string
	pub    *kafka.Writer
	sub    *kafka.Reader
}

func (s *DbMessageCreated) writer() *kafka.Writer {
	if s.pub != nil {
		return s.pub
	}

	cfg := kafka.WriterConfig{
		Brokers: []string{s.broker},
		Topic:   "db.message.created",
	}
	w := kafka.NewWriter(cfg)
	//todo: fails on the first message publish
	w.AllowAutoTopicCreation = true

	s.pub = w
	return w
}

func NewDbMessageCreated(brokerUrl string) *DbMessageCreated {
	return &DbMessageCreated{
		broker: brokerUrl,
	}
}

func (s *DbMessageCreated) Publish(ctx context.Context, value uuid.UUID) error {
	w := s.writer()
	km := kafka.Message{Key: nil, Value: value.Bytes()}
	err := w.WriteMessages(ctx, km)
	return err
}

// func (s *DbMessageCreated) Read()
