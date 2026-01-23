package tgbot

import (
	"gopkg.in/telebot.v3"
	"time"
)

type TgbotClient struct {
	*telebot.Bot
}

func NewTgbotClient(token string) (*TgbotClient, error) {

	b, err := telebot.NewBot(telebot.Settings{
		Token:   token,
		Verbose: true,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return &TgbotClient{
		Bot: b,
	}, nil
}

type SendMessageParams struct {
	ToUserID int64
	Message  Message
}

type Message struct {
	Content string
	Mode    string // Markdown MarkdownV2 HTML
}

func (t *TgbotClient) SendMessage(params SendMessageParams) (*telebot.Message, error) {
	return t.Send(
		&telebot.User{ID: params.ToUserID},
		params.Message.Content,
		params.Message.Mode,
	)
}
