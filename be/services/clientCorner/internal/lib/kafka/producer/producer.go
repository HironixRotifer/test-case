package producer

import (
	"context"
	"encoding/json"
	"log"

	"clientCorner/internal/config"
	"clientCorner/internal/lib/kafka/message"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	addr []string
	conn *kafka.Conn
}

func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	// to produce messages
	topic := "client-corner"
	partition := 1

	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.KafkaAddrs[0], topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return &KafkaProducer{conn: conn}
}

func (k *KafkaProducer) CloseConnKafka() error {
	return k.conn.Close()
}

// SendMessageToKafka send message to kafka
// TODO: make batcher/msg...
func (k *KafkaProducer) SendMessageToKafka(key string, msg message.Message) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(k.addr...),
		Topic:    "client-corner",
		Balancer: &kafka.LeastBytes{},
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: b,
		},
	)

	if err != nil {
		log.Fatal("failed to write messages:", err)
		return err
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
		return err
	}

	return nil
}
