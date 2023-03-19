package sender

import (
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/message"
	"github.com/klzwii/mirai-go/record"
	log "github.com/sirupsen/logrus"
)

func GetWSSender(conn function.Conn, sessionKey string) Sender {
	return &senderWSImp{
		sessionKey: sessionKey,
		conn:       conn,
	}
}

func (d *senderWSImp) SendToFriend(target uint64, contents *message.Chain, quote *uint64) (*record.SendMessageResponseData, error) {
	resp := record.GetSendMessageResponse()
	err := d.conn.SendRequest("sendFriendMessage", "", SendRequest{
		Target:       target,
		MessageChain: contents,
		Quote:        quote,
	}, resp)
	log.Debugf("send to friend get response %+v", resp.GetData().(*record.SendMessageResponseData))
	return resp.GetData().(*record.SendMessageResponseData), err
}

func (d *senderWSImp) SendToGroup(target uint64, contents *message.Chain, quote *uint64) (*record.SendMessageResponseData, error) {
	resp := record.GetSendMessageResponse()
	err := d.conn.SendRequest("sendGroupMessage", "", SendRequest{
		Target:       target,
		MessageChain: contents,
		Quote:        quote,
	}, resp)
	log.Debugf("send to group get response %+v", resp.GetData().(*record.SendMessageResponseData))
	return resp.GetData().(*record.SendMessageResponseData), err
}

type SendRequest struct {
	SessionKey   string         `json:"sessionKey,omitempty"`
	Target       uint64         `json:"target"`
	MessageChain *message.Chain `json:"messageChain"`
	Quote        *uint64        `json:"quote,omitempty"`
}
