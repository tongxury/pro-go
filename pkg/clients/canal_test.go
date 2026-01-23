package clients

import (
	"gopkg.in/telebot.v4"
	"testing"
)

func TestNewCanalClient(t *testing.T) {

	var token = "7221456431:AAFVBeYyYwLM6b_88NbXgm3Spy3fudP9cqw"

	b, _ := telebot.NewBot(telebot.Settings{
		Token: token,
		//Verbose: true,
		//Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	b.Send(&telebot.User{ID: 5459275634}, "aaa")
}

func TestNewCanalClient2(t *testing.B) {

	var token = "7221456431:AAFVBeYyYwLM6b_88NbXgm3Spy3fudP9cqw"

	b, _ := telebot.NewBot(telebot.Settings{
		Token: token,
		//Verbose: true,
		//Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	b.Send(&telebot.User{ID: 5459275634}, "aaa")
}
