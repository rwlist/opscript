package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rwlist/opscript/models"
	"github.com/rwlist/opscript/more"
	"log"
)

type Configure struct {
	chs     *Chats
	repl    *Reply
	dialogs map[int]*models.Dialog
}

func NewConfigure(chs *Chats, repl *Reply) *Configure {
	return &Configure{
		chs:     chs,
		repl:    repl,
		dialogs: make(map[int]*models.Dialog),
	}
}

func (c *Configure) getDialog(id int) (*models.Dialog, error) {
	dialog, ok := c.dialogs[id]
	if !ok {
		dialog = &models.Dialog{}
		c.dialogs[id] = dialog
	}

	return dialog, nil
}

func (c *Configure) HandlePrivate(msg *tgbotapi.Message) {
	dialog, err := c.getDialog(msg.From.ID)
	if err != nil {
		log.Println(err)
		return
	}

	if dialog.Ns == nil {
		c.repl.SendText(msg.Chat.ID, "Write /config in chat you want to configure")
		return
	}

	ns := dialog.Ns
	opts := more.ParseOpts(msg)

	if dialog.Status == models.WaitSrc {
		err := ns.Eval(msg.Text)
		if err != nil {
			c.repl.SendText(msg.Chat.ID, err.Error())
			return
		}
		c.repl.SendText(msg.Chat.ID, "Eval OK")
		dialog.Status = models.Std
		return
	}

	if dialog.Status != models.Std {
		resp := fmt.Sprintf("Wrong status: %v", dialog.Status)
		c.repl.SendText(msg.Chat.ID, resp)
		return
	}

	help := `
Help:

/eval
/exit`

	if !opts.IsCommand {
		c.repl.SendText(msg.Chat.ID, help)
		return
	}

	if opts.Args[0] == "eval" {
		dialog.Status = models.WaitSrc
		c.repl.SendText(msg.Chat.ID, "Enter code to eval")
		return
	}

	if opts.Args[0] == "exit" {
		dialog = &models.Dialog{}
		return
	}

	c.repl.SendText(msg.Chat.ID, help)
	return
}

func (c *Configure) SwitchDialog(user *tgbotapi.User, ns *models.Namespace) {
	dialog, err := c.getDialog(user.ID)
	if err != nil {
		log.Println(err)
		return
	}

	dialog.Ns = ns

	info := fmt.Sprintf("You are now configuring %s", ns.GetName())
	c.repl.SendText(int64(user.ID), info)
}
