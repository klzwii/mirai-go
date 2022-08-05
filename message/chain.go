package message

type Chain []BaseInterface

func (c *Chain) addPlain(text string) *Chain {
	ret := append(*c, &PlainMessage{Text: text})
	return &ret
}
