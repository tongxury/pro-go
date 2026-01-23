package clients

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func (t *KafkaClient) W() *Writer {
	// Make a writer that publishes messages to topic-A.
	// The topic will be created if it is missing.

	w := &kafka.Writer{
		Addr: kafka.TCP(t.brokers...),
		//Topic:                  "topic-A",
		AllowAutoTopicCreation: true,
		Transport: &kafka.Transport{
			SASL: t.mechanism,
		},
	}

	return &Writer{
		writer: w,
	}
}

type Writer struct {
	writer *kafka.Writer
}

func (t *Writer) Write(ctx context.Context, topic string, messages ...kafka.Message) error {

	var msgs []kafka.Message
	for _, x := range messages {
		msgs = append(msgs, kafka.Message{
			Topic:         topic,
			Partition:     x.Partition,
			Offset:        x.Offset,
			HighWaterMark: x.HighWaterMark,
			Key:           x.Key,
			Value:         x.Value,
			Headers:       x.Headers,
			WriterData:    x.WriterData,
			Time:          x.Time,
		})
	}

	return t.writer.WriteMessages(ctx, msgs...)
}
