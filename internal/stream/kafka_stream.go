package stream

import (
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type kafkaStream struct {
	broker   string
	topic    string
	maxBytes int
	pub      *kafka.Writer
	sub      *kafka.Reader
}

func newStream(brokerUrl string, topic string, messageSize int) *kafkaStream {
	return &kafkaStream{
		broker:   brokerUrl,
		topic:    topic,
		maxBytes: messageSize,
	}
}

func (s *kafkaStream) close() error {
	// it is more important to close the reader first
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

func (s *kafkaStream) writer() *kafka.Writer {
	if s.pub != nil {
		return s.pub
	}

	cfg := kafka.WriterConfig{
		Brokers: []string{s.broker},
		Topic:   s.topic,
	}
	w := kafka.NewWriter(cfg)

	s.pub = w
	return w
}

func (s *kafkaStream) reader() *kafka.Reader {
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

func (s *kafkaStream) createTopic() error {
	log.Println("kafkaUrl", s.broker)
	conn, err := kafka.Dial("tcp", s.broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}

	for _, pt := range partitions {
		if pt.Topic == s.topic {
			return nil
		}
	}

	broker, err := conn.Controller()
	if err != nil {
		return err
	}
	brokerUrl := net.JoinHostPort(broker.Host, strconv.Itoa(broker.Port))
	log.Println("kafkaUrl", brokerUrl)

	conn, err = kafka.Dial("tcp", brokerUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	config := kafka.TopicConfig{
		Topic:             s.topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = conn.CreateTopics(config)
	return err
}
