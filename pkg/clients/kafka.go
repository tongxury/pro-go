package clients

import (
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"store/pkg/confcenter"
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
		mechanism, err := scram.Mechanism(scram.SHA512, conf.Username, conf.Password)
		if err != nil {
			panic(err)
		}

		c.mechanism = mechanism
	}

	return c
}
