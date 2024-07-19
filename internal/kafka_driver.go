package internal

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/gofrs/uuid"
	"github.com/lovoo/goka"
)

type KafkaDriver struct {
	emitter *goka.Emitter
}

func NewKafkaDriver(topic string, url string) (*KafkaDriver, error) {
	cfg := goka.DefaultConfig()
	cfg.Version = sarama.V3_5_0_0
	goka.ReplaceGlobalConfig(cfg)

	brokers := []string{url}

	emitter, err := goka.NewEmitter(brokers, goka.Stream(topic), new(UuidCodec))
	if err != nil {
		return nil, fmt.Errorf("ERROR: could not create a kafka/goka emitter: %v", err)
	}

	return &KafkaDriver{
		emitter: emitter,
	}, nil
}

// func (d *KafkaDriver) EmitNewMessage(msg Message) error {
// 	return d.emitId("newMessage", msg.Id)
// }

func (d *KafkaDriver) EmitId(key string, id uuid.UUID) error {
	err := d.emitter.EmitSync(key, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *KafkaDriver) Cleanup() {
	d.emitter.Finish()
}
