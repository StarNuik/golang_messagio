package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/segmentio/kafka-go"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/message"
	"github.com/starnuik/golang_messagio/internal/model"
)

var workloads *model.WorkloadsModel
var messages *model.MessagesModel

func work(km kafka.Message) {
	id := uuid.FromBytesOrNil(km.Value)

	msg, err := messages.Get(context.TODO(), id)
	cmd.ServerError(err)

	load, err := message.Process(msg)
	cmd.ServerError(err)

	err = workloads.Insert(context.TODO(), load)
	cmd.ServerError(err)

	fmt.Printf("received %s, message is: %v\n", km.Topic, msg)
	// fmt.Printf("message received, topic %s, offset %d: %s = %s\n", km.Topic, km.Offset, km.Key, id.String())
}

func main() {
	db, err := internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ServerError(err)
	defer db.Close()

	workloads = model.NewWorkloadsModel(db)
	messages = model.NewMessagesModel(db)

	cfg := kafka.ReaderConfig{
		Brokers: []string{os.Getenv("SERVICE_KAFKA_URL")},
		Topic:   "db.message.created",
		// MaxBytes: 10e3,
		MaxBytes: 16 * 4,
	}
	r := kafka.NewReader(cfg)
	defer r.Close()

	for {
		km, err := r.ReadMessage(context.Background())
		cmd.ServerError(err)
		work(km)
	}
}
