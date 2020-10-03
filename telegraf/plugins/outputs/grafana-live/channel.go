package grafanalive

// Simple Go chat client for https://github.com/centrifugal/centrifuge/tree/master/examples/events example.

import (
	"log"
	"time"

	"github.com/centrifugal/centrifuge-go"
)

// GrafanaLiveChannel lets you write data
type GrafanaLiveChannel struct {
	Channel  string
	sub      *centrifuge.Subscription
	lastWarn time.Time
}

// Publish sends the data to the channel
func (b *GrafanaLiveChannel) Publish(data []byte) {
	status := b.sub.Status()
	if status == centrifuge.SUBSCRIBED {
		b.sub.Publish(data)
	} else if time.Since(b.lastWarn) > time.Second*5 {
		b.lastWarn = time.Now()
		log.Printf("grafana live channel not connected: %s\n", b.Channel)
	}
}

type eventHandler struct {
}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	log.Printf("Connected to chat with ID %s", e.ClientID)
	return
}

func (h *eventHandler) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	log.Printf("Error: %s", e.Message)
	return
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Printf("Disconnected from chat: %s", e.Reason)
	return
}

func (h *eventHandler) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	log.Printf("Subscribed on channel %s, resubscribed: %v, recovered: %v", sub.Channel(), e.Resubscribed, e.Recovered)
}

func (h *eventHandler) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	log.Printf("Subscribed on channel %s failed, error: %s", sub.Channel(), e.Error)
}

func (h *eventHandler) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	log.Printf("Unsubscribed from channel %s", sub.Channel())
}

// InitGrafanaLiveChannel starts a chat server
func InitGrafanaLiveChannel(url string, channel string) (*GrafanaLiveChannel, error) {
	log.Printf("Connect to %s\n", url)

	b := &GrafanaLiveChannel{
		Channel: channel,
	}

	c := centrifuge.New(url, centrifuge.DefaultConfig())
	handler := &eventHandler{}
	c.OnConnect(handler)
	c.OnError(handler)
	c.OnDisconnect(handler)

	sub, err := c.NewSubscription(channel)
	if err != nil {
		return nil, err
	}
	b.sub = sub

	sub.OnSubscribeSuccess(handler)
	sub.OnSubscribeError(handler)
	sub.OnUnsubscribe(handler)

	err = sub.Subscribe()
	if err != nil {
		return nil, err
	}

	err = c.Connect()
	if err != nil {
		return nil, err
	}

	return b, nil
}
