package stream

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/starnuik/golang_messagio/internal/model"
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
	r.SetOffset(kafka.LastOffset)

	s.sub = r
	return r
}

func NewDbMessageCreated(brokerUrl string, messageSize int) (*DbMessageCreated, error) {
	const topic = "db.message.created"

	conn, err := kafka.Dial("tcp", brokerUrl)
	if err != nil {
		return nil, err
	}
	conn.Close()

	return &DbMessageCreated{
		broker:   brokerUrl,
		maxBytes: messageSize,
		topic:    topic,
	}, nil
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
