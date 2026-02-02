package clients

import (
	"store/pkg/confcenter"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type KafkaClient struct {
	mechanism sasl.Mechanism
	brokers   []string
}

func NewKafkaClient(conf confcenter.KafkaConfig) *KafkaClient {

	c := &KafkaClient{
		brokers: conf.Brokers,
	}
	if conf.Password != "" {
		c.mechanism = plain.Mechanism{
			Username: conf.Username,
			Password: conf.Password,
		}
	}

	return c
}
