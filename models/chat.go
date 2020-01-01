package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Chat struct {
	desc *tgbotapi.Chat
	nss  []*Namespace
}

func NewChat(desc *tgbotapi.Chat) *Chat {
	return &Chat{
		desc: desc,
	}
}

func (c *Chat) IsPrivate() bool {
	return c.desc.Type == "private"
}

func (c *Chat) IsGroup() bool {
	return c.desc.Type == "group" || c.desc.Type == "supergroup"
}

func (c *Chat) FindNamespaceByName(name string) *Namespace {
	for _, ns := range c.nss {
		if ns.Name == name {
			return ns
		}
	}

	return nil
}

func (c *Chat) AddNamespace(ns *Namespace) {
	c.nss = append(c.nss, ns)
}

func (c *Chat) AllNamespaces() []*Namespace {
	return c.nss
}

func (c *Chat) GetID() int64 {
	return c.desc.ID
}
