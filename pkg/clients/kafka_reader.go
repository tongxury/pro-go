package clients

import (
	"github.com/segmentio/kafka-go"
	"time"
)

func (t *KafkaClient) R(group, topic string) *kafka.Reader {

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: t.mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  t.brokers,
		GroupID:  group,
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
		Dialer:   dialer,
	})

	return r
}
