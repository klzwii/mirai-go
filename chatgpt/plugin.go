package chatgpt

import (
	"github.com/klzwii/mirai-go/message"
	"github.com/klzwii/mirai-go/sender"
	"github.com/mitchellh/mapstructure"
	"strings"
	"sync"
)
import record2 "github.com/klzwii/mirai-go/record"

type Plugin struct {
	sender sender.Sender
	g      *gpt
}

func (c *Plugin) queryGPT(key string, message string, promptID uint64) string {
	rec := &record{
		Role:    "user",
		Content: message,
	}
	req := c.g.createRequest(rec, key, promptID)
	resp := &response{}
	if err := c.g.sendRequest(req, resp); err != nil {
		logger.Errorf("网络请求失败%v", err)
		return internalError.Error()
	}
	ret, err := c.g.processResponse(resp, rec, key)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(ret.Content)
}

func (c *Plugin) OnGroupMessage(record *record2.GroupMessageData) {
	chain := record.MessageChain
	if len(chain) != 3 || chain[1].GetType() != message.AT || chain[1].(*message.AtMessage).Target != 2844255154 ||
		chain[2].GetType() != message.PLAIN {
		return
	}
	plainMessage := chain[2].(*message.PlainMessage)
	key := generateKey(record.Sender.MemberName, record.Sender.Group.ID)
	ret := c.queryGPT(key, plainMessage.Text, record.Sender.Group.ID)
	_, _ = c.sender.SendToGroup(record.Sender.Group.ID, message.NewMessageChain().AddPlain(ret), &(record.MessageChain[0].(*message.SourceMessage).ID))
}

func (c *Plugin) OnFriendMessage(record *record2.FriendMessageData) {
	chain := record.MessageChain
	if len(chain) != 2 || chain[1].GetType() != message.PLAIN || record.Sender.ID != 1027898733 {
		return
	}
	plainMessage := chain[1].(*message.PlainMessage)
	if plainMessage.Text == "erase user history" {
		logger.Info("用户指令清除全部历史数据")
		c.g.clearHistory()
		logger.Info("清除全部历史数据结束")
		return
	}
	key := generateKey("", record.Sender.ID)
	ret := c.queryGPT(key, plainMessage.Text, record.Sender.ID)
	_, _ = c.sender.SendToFriend(record.Sender.ID, message.NewMessageChain().AddPlain(ret), nil)
}

func (c *Plugin) Name() string {
	return "chatgpt"
}

func (c *Plugin) Initiate(rawConf map[string]any, sender sender.Sender) {
	conf := &config{}
	var (
		decoder *mapstructure.Decoder
		err     error
	)
	if decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{ErrorUnset: true, WeaklyTypedInput: true, Result: conf}); err != nil {
		panic(err)
	}
	if err = decoder.Decode(rawConf); err != nil {
		panic(err)
	}
	c.sender = sender
	c.g = &gpt{
		conf:        conf,
		UserHistory: make(map[string][]*record),
		m:           &sync.Mutex{},
	}
}
