package clients

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"store/pkg/sdk/helper"
)

type Consumer struct {
	kc       *KafkaClient
	handlers []ConsumerHandler
}

type ConsumerHandler struct {
	Group, Topic string
	Handle       func(message kafka.Message) error
}

func (t *Consumer) Consume(handler ConsumerHandler) *Consumer {
	t.handlers = append(t.handlers, handler)
	return t
}

func (t *Consumer) Run() {
	for i := range t.handlers {

		x := t.handlers[i]

		log.Debugw("consumer running", "", "g", x.Group, "t", x.Topic)

		go func(handler ConsumerHandler) {
			defer helper.DeferFunc()

			r := t.kc.R(handler.Group, handler.Topic)

			for {

				err := helper.TryCatch(func() {

					m, err := r.ReadMessage(context.Background())
					if err != nil {
						log.Errorw("failed to ReadMessage:", err, "config", t.kc, "group", handler.Group, "topic", handler.Topic)
						return
					}

					//log.Debugf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

					err = handler.Handle(m)
					if err != nil {
						log.Errorw("Consumer handle err", err, "group", handler.Topic, "topic", handler.Topic, "m", m)
					}

				})

				if err != nil {
					log.Errorw("consumer err", err)
				}
			}
		}(x)
	}

	select {}
}

func (t *KafkaClient) C() *Consumer {
	return &Consumer{
		kc: t,
	}
}
