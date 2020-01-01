package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rwlist/opscript/more"
	"log"
)

type Root struct {
	chs *Chats // long live C.H.S.
	cfg *Configure
}

func NewRoot(chs *Chats, cfg *Configure) *Root {
	return &Root{
		chs: chs,
		cfg: cfg,
	}
}

func (r *Root) Handle(msg *tgbotapi.Message) {
	text := msg.Text
	if text == "" {
		// skip
		return
	}

	spew.Dump(msg.Text, msg.From.UserName)

	chat, err := r.chs.Get(msg.Chat)
	if err != nil {
		log.Println(err)
		return
	}

	opts := more.ParseOpts(msg)
	if opts.IsCommand && opts.Args[0] == "config" && chat.IsGroup() {
		nsName := fmt.Sprintf("ns%v", msg.From.ID)
		ns, err := r.chs.LoadNamespace(chat, nsName, msg.From.ID)
		if err != nil {
			log.Println(err)
			return
		}

		err = r.cfg.SwitchDialog(msg.From, ns)
		if err != nil {
			r.chs.repl.SendText(msg.Chat.ID, "Please write /start to this bot PM to receive messages.")
			return
		}

		return
	}

	if chat.IsPrivate() {
		r.cfg.HandlePrivate(msg)
		return
	}

	r.chs.ExecuteAll(chat, msg)
}
