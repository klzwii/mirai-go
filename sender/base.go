package sender

import (
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/message"
	"github.com/klzwii/mirai-go/record"
)

type Sender interface {
	SendToFriend(target uint64, contents *message.Chain) (*record.SendMessageResponseData, error)
	SendToGroup(target uint64, contents *message.Chain) (*record.SendMessageResponseData, error)
}

type senderWSImp struct {
	sessionKey string
	conn       function.Conn
}
