package message

type Chain []Base

func NewMessageChain() *Chain {
	return &Chain{}
}

func (c *Chain) AddPlain(text string) *Chain {
	ret := append(*c, &PlainMessage{BaseImp: BaseImp{Type: PLAIN}, Text: text})
	return &ret
}
