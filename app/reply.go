package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Reply struct {
	bot *tgbotapi.BotAPI
}

func NewReply(bot *tgbotapi.BotAPI) *Reply {
	return &Reply{bot: bot}
}

func (r *Reply) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := r.bot.Send(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
