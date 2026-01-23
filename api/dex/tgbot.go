package dexpb

import (
	"gopkg.in/telebot.v3"
	"store/pkg/sdk/conv"
)

func (t *BotContext) UserId() string {
	return t.GetUser().GetId()
}

func FromTelebotContext(c telebot.Context) *BotContext {

	sender := c.Sender()

	message := c.Message()

	return &BotContext{
		User: &BotUser{
			Id:           conv.String(sender.ID),
			Firstname:    sender.FirstName,
			Lastname:     sender.LastName,
			LanguageCode: sender.LanguageCode,
		},
		Message: &BotMessage{
			Text:    message.Text,
			Payload: message.Payload,
		},
	}
}

func (t *BotResponse) AsSendOptions() *telebot.SendOptions {

	if t == nil {
		return nil
	}

	return &telebot.SendOptions{
		ReplyTo:               nil,
		ReplyMarkup:           nil,
		DisableWebPagePreview: false,
		DisableNotification:   false,
		ParseMode:             "",
		Entities:              nil,
		AllowWithoutReply:     false,
		Protected:             false,
		ThreadID:              0,
		HasSpoiler:            false,
		ReplyParams:           nil,
		//BusinessConnectionID:  "",
		//EffectID:              "",
	}
}

// Markets
//type Markets []*Market

//func (ts Markets) Tokens() Tokens {
//	var tokens Tokens
//	for _, t := range ts {
//
//		t.Token.MarketState = &MarketState{
//			Stats5M:  nil,
//			Stats1H:  nil,
//			Stats6H:  nil,
//			Stats24H: nil,
//		}
//
//		tokens = append(tokens, t.Token)
//	}
//	return tokens
//}

// UserWalletRelations
type UserWalletRelations []*UserWalletRelation

func (ts UserWalletRelations) Wallets() Wallets {
	var wallets Wallets
	for _, wallet := range ts {
		wallets = append(wallets, wallet.Wallet)
	}

	return wallets
}

func (ts UserWalletRelations) UserIds() []string {
	var wallets []string
	for _, wallet := range ts {
		wallets = append(wallets, wallet.User.XId)
	}

	return wallets
}

// Amounts
type Amounts []*Amount

func (ts Amounts) AsMap() map[string]*Amount {
	result := make(map[string]*Amount, len(ts))
	for _, t := range ts {
		result[t.XId] = t
	}
	return result
}

// Token
type Tokens []*Token

func (ts Tokens) Prices() Amounts {

	var tokens Amounts

	for _, x := range ts {
		tokens = append(tokens, x.Price)
	}

	return tokens
}

func (ts Tokens) TokenIds() []string {

	var tokens []string

	for _, x := range ts {
		tokens = append(tokens, x.XId)
	}

	return tokens
}

func (ts Tokens) AsMap() map[string]*Token {
	tokens := make(map[string]*Token, len(ts))
	for _, t := range ts {
		tokens[t.XId] = t
	}
	return tokens
}

// UserTokenRelations
type UserTokenRelations []*UserTokenRelation

func (ts UserTokenRelations) AsMap() map[string]*UserTokenRelation {
	result := make(map[string]*UserTokenRelation, len(ts))
	for _, t := range ts {
		result[t.User.XId+t.Token.XId] = t
	}
	return result
}

func (ts UserTokenRelations) TokenIds() []string {

	var tokens []string

	for _, x := range ts {
		tokens = append(tokens, x.Token.XId)
	}

	return tokens
}
