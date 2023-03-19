package bot

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/sender"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"net/url"
)

var logger = logrus.WithField("core", "bot")

type Bot struct {
	sender  sender.Sender
	conn    function.Conn
	plugins []Plugin
	conf    map[string]interface{}
}

type config struct {
	Port      string
	VerifyKey string `mapstructure:"verify_key"`
	QQ        string
}

func (b *Bot) Start(ctx context.Context) {
	initiatePlugins(b.conf, b.plugins, b.sender)
	ch := make(chan record.Base, 100)
	go b.conn.StartReading(ctx, ch)
	logger.Info("bot初始化完成")
	for {
		select {
		case curRecord := <-ch:
			go triggerPlugins(curRecord, b)
		case <-ctx.Done():
			logger.Info("bot退出")
			return
		}
	}
}

func (b *Bot) RegisterPlugin(plugin Plugin) {
	b.plugins = append(b.plugins, plugin)
}

func GetBot(rawConfig map[string]any) (*Bot, error) {
	var (
		conf    = &config{}
		err     error
		decoder *mapstructure.Decoder
	)
	if decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnset: true, WeaklyTypedInput: true, Result: conf}); err != nil {
		return nil, err
	}
	if err = decoder.Decode(rawConfig["bot"]); err != nil {
		return nil, err
	}
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("localhost:%s", conf.Port), Path: "/message"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), map[string][]string{
		"verifyKey": {conf.VerifyKey},
		"qq":        {conf.QQ},
	})
	if err != nil {
		return nil, err
	}
	conn := function.GetWsConn(c)
	return &Bot{
		sender.GetWSSender(conn, ""),
		conn,
		[]Plugin{},
		rawConfig,
	}, nil
}
