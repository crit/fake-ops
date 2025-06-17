package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Context struct {
	context.Context
	publish func(msg tea.Msg)

	Flags Flags
}

func NewContext(sender func(msg tea.Msg)) (*Context, func()) {
	var flags Flags
	flags.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	return &Context{
		Context: ctx,
		publish: sender,
		Flags:   flags,
	}, cancel
}

func (ctx *Context) PublishInfo(msg string, args ...any) {
	ctx.publish(Message{
		Kind:  InfoKind,
		Value: fmt.Sprintf(msg, args...)},
	)
}

func (ctx *Context) PublishError(msg string, args ...any) {
	ctx.publish(Message{
		Kind:  ErrorKind,
		Value: fmt.Sprintf(msg, args...)},
	)
}

func (ctx *Context) PublishService(kind, name string, port int) {
	ctx.publish(ServiceMessage{
		Kind:       kind,
		Name:       name,
		Port:       port,
		Status:     "offline",
		LastStatus: time.Now(),
	})
}

func (ctx *Context) PublishServiceOnline(name string) {
	ctx.publish(ServiceStatus{
		Sent:   time.Now(),
		Name:   name,
		Status: "online",
	})
}

func (ctx *Context) PublishServiceOffline(name string) {
	ctx.publish(ServiceStatus{
		Sent:   time.Now(),
		Name:   name,
		Status: "offline",
	})
}

func (ctx *Context) PublishServiceError(name string) {
	ctx.publish(ServiceStatus{
		Sent:   time.Now(),
		Name:   name,
		Status: "error",
	})
}
