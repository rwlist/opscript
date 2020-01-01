package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rwlist/opscript/models"
	"github.com/rwlist/opscript/more"
)

type Chats struct {
	repl  *Reply
	chats map[int64]*models.Chat
}

func NewChats(repl *Reply) *Chats {
	return &Chats{
		repl:  repl,
		chats: make(map[int64]*models.Chat),
	}
}

func (c *Chats) Create(desc *tgbotapi.Chat) (*models.Chat, error) {
	// TODO: fetch from db
	return models.NewChat(desc), nil
}

func (c *Chats) Get(desc *tgbotapi.Chat) (*models.Chat, error) {
	chat, ok := c.chats[desc.ID]
	if !ok {
		var err error
		chat, err = c.Create(desc)
		if err != nil {
			return nil, err
		}

		c.chats[desc.ID] = chat
	}

	return chat, nil
}

func (c *Chats) LoadNamespace(chat *models.Chat, name string) (*models.Namespace, error) {
	ns := chat.FindNamespaceByName(name)
	if ns == nil {
		ns = models.NewNamespace(name, chat.GetID())
		chat.AddNamespace(ns)
	}

	return ns, nil
}

func (c *Chats) ExecuteAll(chat *models.Chat, msg *tgbotapi.Message) {
	opts := more.ParseOpts(msg)
	if !opts.IsCommand {
		return
	}

	nss := chat.AllNamespaces()
	for _, ns := range nss {
		resp, err := ns.Act(opts.Args)
		if err != nil {
			c.repl.SendText(chat.GetID(), err.Error())
			continue
		}

		if resp == "" {
			continue
		}

		c.repl.SendText(chat.GetID(), resp)
	}
}
