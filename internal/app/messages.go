package app

import "time"

type MessageKind string

const (
	ErrorKind MessageKind = "error"
	InfoKind  MessageKind = "info"
)

// Message communicates log information.
type Message struct {
	Value string
	Kind  MessageKind
}

func (msg Message) String() string {
	return msg.Value
}

// ServiceMessage communicates about a service.
type ServiceMessage struct {
	Kind       string
	Name       string
	Port       int
	Status     string
	LastStatus time.Time
}

// ServiceStatus communicates about a service's status.
type ServiceStatus struct {
	Sent   time.Time
	Name   string
	Status string
}
