package bot

import (
	"fmt"
	"github.com/klzwii/mirai-go/record"
)

type Plugin interface {
	OnGroupMessage(record *record.GroupMessageData)
	OnFriendMessage(record *record.FriendMessageData)
}

type TestPlugin struct {
}

func (t *TestPlugin) OnGroupMessage(record *record.GroupMessageData) {
	fmt.Println("this is a group Message", record.Sender)
}

func (t *TestPlugin) OnFriendMessage(record *record.FriendMessageData) {
	fmt.Println("this is a friend Message", record.Sender)
}
