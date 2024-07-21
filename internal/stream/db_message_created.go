package stream

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/segmentio/kafka-go"
)

type DbMessageCreated struct {
	broker   string
	topic    string
	maxBytes int
	pub      *kafka.Writer
	sub      *kafka.Reader
}

func (s *DbMessageCreated) writer() *kafka.Writer {
	if s.pub != nil {
		return s.pub
	}

	cfg := kafka.WriterConfig{
		Brokers: []string{s.broker},
		Topic:   s.topic,
	}
	w := kafka.NewWriter(cfg)
	w.AllowAutoTopicCreation = true

	s.pub = w
	return w
}

func (s *DbMessageCreated) reader() *kafka.Reader {
	if s.sub != nil {
		return s.sub
	}

	cfg := kafka.ReaderConfig{
		Brokers:  []string{s.broker},
		Topic:    s.topic,
		MaxBytes: s.maxBytes,
	}
	r := kafka.NewReader(cfg)

	s.sub = r
	return r
}

func NewDbMessageCreated(brokerUrl string, messageSize int) *DbMessageCreated {
	return &DbMessageCreated{
		broker:   brokerUrl,
		maxBytes: messageSize,
		topic:    "db.message.created",
	}
}

func (s *DbMessageCreated) Publish(ctx context.Context, value uuid.UUID) error {
	w := s.writer()
	km := kafka.Message{Key: nil, Value: value.Bytes()}
	err := w.WriteMessages(ctx, km)
	return err
}

func (s *DbMessageCreated) Read(ctx context.Context) (uuid.UUID, error) {
	r := s.reader()
	km, err := r.ReadMessage(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.FromBytes(km.Value)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *DbMessageCreated) Close() error {
	// its more important to close the reader first
	// to allow the broker to send messages to other consumers
	if s.sub != nil {
		err := s.sub.Close()
		if err != nil {
			return err
		}
		s.sub = nil
	}

	if s.pub != nil {
		err := s.pub.Close()
		if err != nil {
			return err
		}
		s.pub = nil
	}
	return nil
}

func (s *DbMessageCreated) Topic() string {
	return s.topic
}
