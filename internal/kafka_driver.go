package internal

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type KafkaDriver struct {
	emitter *goka.Emitter
}

func NewKafkaDriver(url string) (*KafkaDriver, error) {
	cfg := goka.DefaultConfig()
	cfg.Version = sarama.V3_5_0_0
	goka.ReplaceGlobalConfig(cfg)

	brokers := []string{url}
	topic := goka.Stream("messages")

	emitter, err := goka.NewEmitter(brokers, topic, new(codec.Bytes))
	if err != nil {
		return nil, fmt.Errorf("ERROR: could not create a kafka/goka emitter: %v", err)
	}

	return &KafkaDriver{
		emitter: emitter,
	}, nil
}

func (d *KafkaDriver) Emit(key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ERROR: could not convert object to json: %v %v", value, err)
	}

	err = d.emitter.EmitSync(key, payload)
	if err != nil {
		return fmt.Errorf("ERROR: could not emit a kafka/goka message: %v", err)
	}

	return nil
}

func (d *KafkaDriver) Cleanup() {
	d.emitter.Finish()
}
