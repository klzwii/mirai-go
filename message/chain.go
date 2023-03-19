package message

type Chain []Base

func NewMessageChain() *Chain {
	return &Chain{}
}

func (c *Chain) AddPlain(text string) *Chain {
	ret := append(*c, &PlainMessage{BaseImp: BaseImp{Type: PLAIN}, Text: text})
	return &ret
}

func (c *Chain) AddAt(target uint64, display string) *Chain {
	ret := append(*c, &AtMessage{
		BaseImp: BaseImp{Type: AT},
		Target:  target,
		Display: display,
	})
	return &ret
}

func (c *Chain) AddQuote(id, groupId, senderId, targetId uint64, origin Chain) *Chain {
	ret := append(*c, &QuoteMessage{
		BaseImp:  BaseImp{Type: Quote},
		Id:       id,
		GroupId:  groupId,
		SenderId: senderId,
		TargetId: targetId,
		Origin:   origin,
	})
	return &ret
}
