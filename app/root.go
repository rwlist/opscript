package app

import (
	"fmt"
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

	chat, err := r.chs.Get(msg.Chat)
	if err != nil {
		log.Println(err)
		return
	}

	opts := more.ParseOpts(msg)
	if opts.IsCommand && opts.Args[0] == "config" && chat.IsGroup() {
		nsName := fmt.Sprintf("ns%v", msg.From.ID)
		ns, err := r.chs.LoadNamespace(chat, nsName)
		if err != nil {
			log.Println(err)
			return
		}

		r.cfg.SwitchDialog(msg.From, ns)
		return
	}

	if chat.IsPrivate() {
		r.cfg.HandlePrivate(msg)
		return
	}

	r.chs.ExecuteAll(chat, msg)
}
