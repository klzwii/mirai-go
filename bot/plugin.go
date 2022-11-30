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
	RegisterSender(sender sender.Sender)
}

type TestPlugin struct {
	sender.Sender
}

func (t *TestPlugin) RegisterSender(sender sender.Sender) {
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
	if resp, err := t.Sender.SendToGroup(record.Sender.Group.ID, message.NewMessageChain().AddPlain("testing")); err != nil || resp.Code != 0 {
		log.Errorf("Test plugin send message error %v, resp %v", err, resp)
	}
}

func (t *TestPlugin) OnFriendMessage(record *record.FriendMessageData) {
}
