package sender

import (
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/message"
)

type Sender interface {
	SendToFriend(target uint64, contents *message.Chain) error
	SendToGroup(target uint64, contents *message.Chain) error
}

type senderWSImp struct {
	sessionKey string
	conn       function.Conn
}
