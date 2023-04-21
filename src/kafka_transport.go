package app

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

/***
Kafka transport layer implementation
***/
type KafkaTransport struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewKafkaTransport(incomingTopic string, outgoingTopic string, kafkaAddress string) *KafkaTransport {
	return &KafkaTransport{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{kafkaAddress},
			Topic:     incomingTopic,
			GroupID:   "trade-capture-group",
			Partition: 0,
			MinBytes:  10e3,
			MaxBytes:  10e6,
		}),
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{kafkaAddress},
			Topic:    outgoingTopic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (t *KafkaTransport) Receive(ctx context.Context) (*MessagePayload, error) {
	msg, err := t.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	defer t.reader.CommitMessages(ctx, msg)

	payload := &MessagePayload{
		Key:   string(msg.Key),
		Value: msg.Value,
	}

	return payload, nil
}

func (t *KafkaTransport) Send(ctx context.Context, payload *MessagePayload) error {
	message := kafka.Message{
		Key:   []byte(payload.Key),
		Value: payload.Value,
	}

	err := t.writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	log.Printf("Sent message with key %s", payload.Key)
	return nil
}
