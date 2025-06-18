package app

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Context handles communicating to the UI and processing flags.
type Context struct {
	context.Context
	publish func(msg tea.Msg)

	Flags Flags
}

// NewContext creates a new Context with an attached sender method for communicating
// with the UI.
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

// PublishInfo sends an info Message to the UI.
func (ctx *Context) PublishInfo(msg string, args ...any) {
	ctx.publish(Message{
		Kind:  InfoKind,
		Value: fmt.Sprintf(msg, args...)},
	)
}

// PublishError sends an error Message to the UI.
func (ctx *Context) PublishError(msg string, args ...any) {
	ctx.publish(Message{
		Kind:  ErrorKind,
		Value: fmt.Sprintf(msg, args...)},
	)
}

// PublishService sends a ServiceMessage to the UI. Registering the service
// with the UI.
func (ctx *Context) PublishService(kind, name string, port int) {
	ctx.publish(ServiceMessage{
		Kind:       kind,
		Name:       name,
		Port:       port,
		Status:     "offline",
		LastStatus: time.Now(),
	})
}

// PublishServiceOnline sends a ServiceStatus to the UI indicating that the
// service is online.
func (ctx *Context) PublishServiceOnline(name string) {
	ctx.publish(ServiceStatus{
		Sent:   time.Now(),
		Name:   name,
		Status: "online",
	})
}

// PublishServiceOffline sends a ServiceStatus to the UI indicating that the
// service is offline.
func (ctx *Context) PublishServiceOffline(name string) {
	ctx.publish(ServiceStatus{
		Sent:   time.Now(),
		Name:   name,
		Status: "offline",
	})
}

// PublishServiceError sends a ServiceStatus to the UI indicating that the
// service is currently erroring.
func (ctx *Context) PublishServiceError(name string) {
	ctx.publish(ServiceStatus{
		Sent:   time.Now(),
		Name:   name,
		Status: "error",
	})
}
