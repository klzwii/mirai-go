package bot

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/klzwii/mirai-go/assembler"
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/sender"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type Reader <-chan record.Base

type Bot struct {
	sender.Sender
	Reader
	plugins []Plugin
}

func (b *Bot) Start(ctx context.Context) {
	for _, plugin := range b.plugins {
		plugin.RegisterSender(b.Sender)
	}
	for {
		select {
		case curRecord := <-b.Reader:
			switch curRecord.GetType() {
			case record.FriendMessage:
				message := (curRecord.GetData()).(*record.FriendMessageData)
				for _, plugin := range b.plugins {
					plugin.OnFriendMessage(message)
				}
			case record.GroupMessage:
				message := (curRecord.GetData()).(*record.GroupMessageData)
				for _, plugin := range b.plugins {
					plugin.OnGroupMessage(message)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *Bot) RegisterPlugin(plugin Plugin) {
	b.plugins = append(b.plugins, plugin)
}

func GetBot() (*Bot, error) {
	u := url.URL{Scheme: "ws", Host: "localhost:5679", Path: "/message"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), map[string][]string{
		"verifyKey": {"1234567890"},
		"qq":        {"2844255154"},
	})
	if err != nil {
		return nil, err
	}
	ch := make(chan record.Base, 100)
	go func() {
		for {
			_, m, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Debugf("Get ws message %v", string(m))
			ch <- assembler.UnmarshalToRecord(string(m))
			//log.Printf("recv: %s", m)
		}
	}()
	return &Bot{
		sender.GetWSSender(function.GetWsConn(c), ""),
		ch,
		[]Plugin{},
	}, nil
}
