package models

import (
	"github.com/containous/yaegi/interp"
	"github.com/pkg/errors"
	"strconv"
)

type Namespace struct {
	Name   string
	ChatID int64

	i *interp.Interpreter
}

func NewNamespace(name string, chatID int64) *Namespace {
	ns := &Namespace{
		Name:   name,
		ChatID: chatID,
	}
	ns.Flush()

	return ns
}

func (n *Namespace) GetName() string {
	return strconv.Itoa(int(n.ChatID)) + "::" + n.Name
}

func (n *Namespace) Flush() {
	n.i = interp.New(interp.Options{})
}

func (n *Namespace) Eval(src string) error {
	_, err := n.i.Eval(src)
	return n.wrap(err, "failed to eval")
}

func (n *Namespace) Act(args []string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered from panic %v", r)
		}
	}()

	v, err := n.i.Eval("ns.HandleCmd")
	if err != nil {
		return "", n.wrap(err, "HandleCmd not found")
	}

	f, ok := v.Interface().(func([]string) string)
	if !ok {
		return "", errors.Errorf("cast to (func([]string) string) failed")
	}

	return f(args), nil
}

func (n *Namespace) wrap(err error, comment string) error {
	if err == nil {
		return nil
	}

	return errors.Wrapf(errors.Wrap(err, comment), "error in %s", n.GetName())
}
