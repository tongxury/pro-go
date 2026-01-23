package service

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"store/pkg/clients"
	"store/pkg/events"
)

func (t BiService) StartConsumer() {

	t.data.KafkaClient.C().
		Consume(clients.ConsumerHandler{
			Group: events.TopicPaymentSuccess + "_group",
			Topic: events.TopicPaymentSuccess,
			Handle: func(message kafka.Message) error {

				var event events.PaymentSuccessEvent

				err := json.Unmarshal(message.Value, &event)
				if err != nil {
					return err
				}

				log.Debugw("Consumer ", events.TopicPaymentSuccess, "event", event)

				err = t.eventLog.OnPayEvent(context.Background(), event)
				if err != nil {
					log.Errorw("failed to add event to log , error", err, "event", event)
				}

				return nil
			},
		}).
		Consume(clients.ConsumerHandler{
			Group: events.Topic_AuthLogin + "_group",
			Topic: events.Topic_AuthLogin,
			Handle: func(message kafka.Message) error {

				var event events.AuthEvent

				err := json.Unmarshal(message.Value, &event)
				if err != nil {
					return err
				}

				log.Debugw("Consumer ", events.Topic_AuthLogin, "event", event)

				if event.IsRegister {

					err = t.eventLog.OnRegisterEvent(context.Background(), event)
					if err != nil {
						log.Errorw("failed to add event to log , error", err, "event", event)
					}

				}

				return nil
			},
		}).
		Run()
}
