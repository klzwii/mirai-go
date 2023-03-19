package bot

import (
	"github.com/klzwii/mirai-go/message"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/sender"
	log "github.com/sirupsen/logrus"
)

type Plugin interface {
	OnGroupMessage(record *record.GroupMessageData)
	OnFriendMessage(record *record.FriendMessageData)
	Name() string
	Initiate(conf map[string]any, sender sender.Sender)
}

func initiatePlugins(conf map[string]any, plugins []Plugin, sender sender.Sender) {
	for _, plugin := range plugins {
		var curConf map[string]any = nil
		if conf[plugin.Name()] != nil {
			curConf = conf[plugin.Name()].(map[string]any)
		}
		plugin.Initiate(curConf, sender)
	}
}

func triggerPlugins(curRecord record.Base, b *Bot) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("plugin fatal panic %v", err)
		}
	}()
	switch curRecord.GetType() {
	case record.FriendMessage:
		curMessage := (curRecord.GetData()).(*record.FriendMessageData)
		for _, plugin := range b.plugins {
			plugin.OnFriendMessage(curMessage)
		}
	case record.GroupMessage:
		curMessage := (curRecord.GetData()).(*record.GroupMessageData)
		for _, plugin := range b.plugins {
			plugin.OnGroupMessage(curMessage)
		}
	}
}

type TestPlugin struct {
	sender.Sender
}

func (t *TestPlugin) Initiate(_ map[string]any, sender sender.Sender) {
	t.Sender = sender
}

func (t *TestPlugin) OnGroupMessage(record *record.GroupMessageData) {
	chain := record.MessageChain
	if len(chain) < 2 || chain[1].GetType() != message.PLAIN {
		return
	}
	if chain[1].(*message.PlainMessage).Text != "#抽签" {
		return
	}
	if resp, err := t.Sender.SendToGroup(record.Sender.Group.ID, message.NewMessageChain().AddPlain("testing"), nil); err != nil || resp.Code != 0 {
		log.Errorf("Test plugin send message error %v, resp %v", err, resp)
	}
}

func (t *TestPlugin) Name() string {
	return "TestPlugin"
}

func (t *TestPlugin) OnFriendMessage(_ *record.FriendMessageData) {
}
