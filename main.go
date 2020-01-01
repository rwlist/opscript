package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rwlist/opscript/app"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ltime)

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	repl := app.NewReply(bot)
	chs := app.NewChats(repl)
	cfg := app.NewConfigure(chs, repl)
	handler := app.NewRoot(chs, cfg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		handler.Handle(update.Message)

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID
		//
		// bot.Send(msg)
	}
}
