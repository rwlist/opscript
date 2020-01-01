package models

import (
	"context"
	"github.com/containous/yaegi/interp"
	"github.com/pkg/errors"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type Namespace struct {
	Name    string
	ChatID  int64
	OwnerID int

	i *interp.Interpreter
}

func NewNamespace(name string, chatID int64, ownerID int) *Namespace {
	ns := &Namespace{
		Name:    name,
		ChatID:  chatID,
		OwnerID: ownerID,
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

func (n *Namespace) Eval(src string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered from panic %v.\n%s", r, debug.Stack())
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Second/2)
	_, err = n.i.EvalWithContext(ctx, src)
	return n.wrap(err, "failed to eval")
}

func (n *Namespace) Act(args []string) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered from panic %v.\n%s", r, debug.Stack())
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Second/2)

	v, err := n.i.EvalWithContext(ctx, "ns.HandleCmd")
	if err != nil {
		if strings.Contains(err.Error(), "undefined: ns") {
			return "", nil
		}
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
