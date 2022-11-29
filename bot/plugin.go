package bot

import (
	"fmt"
	"github.com/klzwii/mirai-go/message"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/sender"
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
	fmt.Println("this is a group Message", record.Sender)
}

func (t *TestPlugin) OnFriendMessage(record *record.FriendMessageData) {
	_ = t.Sender.SendToGroup(590258464, message.NewMessageChain().AddPlain("123"))
	fmt.Println("this is a friend Message", record.Sender)
}
