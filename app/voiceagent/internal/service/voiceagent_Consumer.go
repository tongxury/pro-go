package service

import (
	"context"
	"encoding/json"
	"store/pkg/clients"
	"store/pkg/events"
	"store/pkg/sdk/helper"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
)

func (s *VoiceAgentService) StartConsumer() {

	go func() {
		defer helper.DeferFunc()
		s.Data.KafkaClient.C().
			Consume(clients.ConsumerHandler{
				Group: "voiceagent_auth_login_group",
				Topic: events.Topic_AuthLogin,
				Handle: func(message kafka.Message) error {
					var event events.AuthEvent
					err := json.Unmarshal(message.Value, &event)
					if err != nil {
						return err
					}

					if event.IsRegister {
						log.Infow("New user registered, creating default counseling agent", "userId", event.UserID)
						err = s.AgentBiz.CreateDefaultCounselingAgent(context.Background(), event.UserID)
						if err != nil {
							log.Errorw("failed to create default counseling agent", "error", err, "userId", event.UserID)
						}
					}

					return nil
				},
			}).
			Run()
	}()

}
