package api

import (
	"fmt"
	apitypes "lenscape-chat/api-types"

	"github.com/dustin/go-broadcast"
)

type ChannelManager struct {
	channels map[string]broadcast.Broadcaster
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]broadcast.Broadcaster),
	}
}

func (c *ChannelManager) NewSession(channel string) chan any {
	listener := make(chan interface{})
	if _, ok := c.channels[channel]; !ok {
		c.channels[channel] = broadcast.NewBroadcaster(10)
	}
	c.channels[channel].Register(listener)
	return listener
}

func (c *ChannelManager) CreateChannelMessage(channel string, message *apitypes.ChatMessage) {
	if _, ok := c.channels[channel]; !ok {
		c.channels[channel] = broadcast.NewBroadcaster(10)
	}
	fmt.Println("chan", c.channels[channel])
	c.channels[channel].Submit(message)
}
