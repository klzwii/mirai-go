package sender

import (
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/message"
)

func GetWSSender(conn function.Conn, sessionKey string) Sender {
	return &senderWSImp{
		sessionKey: sessionKey,
		conn:       conn,
	}
}

func (d *senderWSImp) SendToFriend(target uint64, contents *message.Chain) error {
	_, err := d.conn.SendRequest("sendFriendMessage", "", SendRequest{
		Target:       target,
		MessageChain: contents,
	})
	return err
}

func (d *senderWSImp) SendToGroup(target uint64, contents *message.Chain) error {
	_, err := d.conn.SendRequest("sendGroupMessage", "", SendRequest{
		Target:       target,
		MessageChain: contents,
	})
	return err
}

type SendRequest struct {
	SessionKey   string         `json:"sessionKey,omitempty"`
	Target       uint64         `json:"target"`
	MessageChain *message.Chain `json:"messageChain"`
}
