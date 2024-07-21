package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/starnuik/golang_messagio/internal"
	"github.com/starnuik/golang_messagio/internal/cmd"
	"github.com/starnuik/golang_messagio/internal/message"
	"github.com/starnuik/golang_messagio/internal/model"
	"github.com/starnuik/golang_messagio/internal/stream"
)

var workloads *model.WorkloadsModel
var messages *model.MessagesModel
var messageCreated *stream.DbMessageCreated

func processMessage(id uuid.UUID) {
	msg, err := messages.Get(context.TODO(), id)
	cmd.ServerError(err)

	load, err := message.Process(msg)
	cmd.ServerError(err)

	err = workloads.Insert(context.TODO(), load)
	cmd.ServerError(err)

	fmt.Printf("processed %s, message was: %v\n", messageCreated.Topic(), msg)
}

func main() {
	db, err := internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ServerError(err)
	defer db.Close()

	workloads = model.NewWorkloadsModel(db)
	messages = model.NewMessagesModel(db)

	messageCreated = stream.NewDbMessageCreated(os.Getenv("SERVICE_KAFKA_URL"), 10e3)
	defer messageCreated.Close()

	for {
		id, err := messageCreated.Read(context.TODO())
		cmd.ServerError(err)
		processMessage(id)
	}
}
