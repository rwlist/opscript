package more

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rwlist/opscript/models"
	"strings"
)

func ParseOpts(msg *tgbotapi.Message) *models.Opts {
	opts := &models.Opts{}

	text := msg.Text
	if s := strings.TrimPrefix(text, "/"); s != text {
		args := strings.Fields(s)

		if len(args) > 0 {
			opts.IsCommand = true
			opts.Args = args
		}
	}

	return opts
}
