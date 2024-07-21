package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/segmentio/kafka-go"
)

func main() {
	cfg := kafka.ReaderConfig{
		Brokers: []string{os.Getenv("SERVICE_KAFKA_URL")},
		Topic:   "db.message.created",
		// MaxBytes: 10e3,
		MaxBytes: 16 * 4,
	}
	r := kafka.NewReader(cfg)
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())

		if err != nil {
			break
		}

		id := uuid.FromBytesOrNil(m.Value)

		fmt.Printf("message received, topic %s, offset %d: %s = %s\n", m.Topic, m.Offset, m.Key, id.String())
	}
}
