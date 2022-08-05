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

func GetWSSender(conn function.Conn, sessionKey string) Sender {
	return &senderWSImp{
		sessionKey: sessionKey,
		conn:       conn,
	}
}

func (d *senderWSImp) SendToFriend(target uint64, contents *message.Chain) error {
	return d.conn.SendRequest("sendFriendMessage", "", SendRequest{
		Target:       target,
		MessageChain: contents,
	})
}

func (d *senderWSImp) SendToGroup(target uint64, contents *message.Chain) error {
	return d.conn.SendRequest("sendGroupMessage", "", SendRequest{
		Target:       target,
		MessageChain: contents,
	})
}

type SendRequest struct {
	SessionKey   string         `json:"sessionKey,omitempty"`
	Target       uint64         `json:"target"`
	MessageChain *message.Chain `json:"messageChain"`
}
