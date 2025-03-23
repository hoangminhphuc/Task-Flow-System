package pubsub

import (
	"context"
	"fmt"
	"time"
)

type Topic string

type PubSub interface {
	Publish(ctx context.Context, channel Topic, data *Message) error
	Subscribe(ctx context.Context, channel Topic) (ch <-chan *Message, close func())
}

type Message struct {
	id        string
	channel   Topic
	data      interface{}
	createdAt time.Time
}

func NewMessage(data interface{}) *Message {
	now := time.Now().UTC()

	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		createdAt: now,
	}
}

func (evt *Message) String() string {
	return fmt.Sprintf("Message %s value %v", evt.channel, evt.data)
}

func (evt *Message) Channel() Topic {
	return evt.channel
}

func (evt *Message) SetChannel(channel Topic) {
	evt.channel = channel
}

func (evt *Message) Data() interface{} {
	return evt.data
}