package app

import "time"

type MessageKind string

const (
	ErrorKind MessageKind = "error"
	InfoKind  MessageKind = "info"
)

type Message struct {
	Value string
	Kind  MessageKind
}

func (msg Message) String() string {
	return msg.Value
}

type ServiceMessage struct {
	Kind       string
	Name       string
	Port       int
	Status     string
	LastStatus time.Time
}

type ServiceStatus struct {
	Sent   time.Time
	Name   string
	Status string
}
