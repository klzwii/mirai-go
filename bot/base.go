package bot

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/sender"
	"net/url"
)

type Bot struct {
	Sender  sender.Sender
	conn    function.Conn
	plugins []Plugin
}

func (b *Bot) Start(ctx context.Context) {
	for _, plugin := range b.plugins {
		plugin.RegisterSender(b.Sender)
	}
	ch := make(chan record.Base, 100)
	go b.conn.StartReading(ctx, ch)
	for {
		select {
		case curRecord := <-ch:
			go func() {
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
			}()
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
	conn := function.GetWsConn(c)
	return &Bot{
		sender.GetWSSender(conn, ""),
		conn,
		[]Plugin{},
	}, nil
}
