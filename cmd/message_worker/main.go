package main

import (
	"context"
	"log"
	"os"
	"time"

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
	start := time.Now()

	isProcessed, err := workloads.Exists(context.TODO(), id)
	cmd.ExitIf(err)
	if isProcessed {
		return
	}

	msg, err := messages.Get(context.TODO(), id)
	cmd.ExitIf(err)

	load, err := message.Process(msg)
	cmd.ExitIf(err)

	err = workloads.Insert(context.TODO(), load)
	cmd.ExitIf(err)

	diff := time.Since(start)
	log.Printf("sucess, duration: %v\n", diff)
}

func main() {
	db, err := internal.NewSqlPool(os.Getenv("SERVICE_POSTGRES_URL"))
	cmd.ExitIf(err)
	defer db.Close()

	workloads = model.NewWorkloadsModel(db)
	messages = model.NewMessagesModel(db)

	messageCreated = stream.NewDbMessageCreated(os.Getenv("SERVICE_KAFKA_URL"), 10e3)
	defer messageCreated.Close()

	for {
		id, err := messageCreated.Read(context.TODO())
		cmd.ExitIf(err)
		// todo: better threading
		go processMessage(id)
	}
}
