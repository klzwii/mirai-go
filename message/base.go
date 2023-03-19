package message

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type Type string

const (
	SOURCE = Type("Source")
	PLAIN  = Type("Plain")
	AT     = Type("At")
	Quote  = Type("Quote")
)

func GetMessageChain(results []gjson.Result) (Chain, error) {
	var chain = Chain{}
	for _, result := range results {
		if ele, err := getMessage(result); err != nil {
			return nil, err
		} else {
			chain = append(chain, ele)
		}
	}
	return chain, nil
}

func getMessage(result gjson.Result) (Base, error) {
	messageType := Type(result.Get("type").Str)
	var ret Base = nil
	switch messageType {
	case SOURCE:
		ret = &SourceMessage{}
	case PLAIN:
		ret = &PlainMessage{}
	case AT:
		ret = &AtMessage{}
	case Quote:
		origin, err := GetMessageChain(result.Get("origin").Array())
		if err != nil {
			return nil, err
		}
		ret = &QuoteMessage{Origin: origin}
	default:
		ret = &BaseImp{}
	}
	if err := jsoniter.UnmarshalFromString(result.Raw, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type Base interface {
	GetType() Type
}

type BaseImp struct {
	Type Type `json:"type,omitempty"`
}

func (b *BaseImp) GetType() Type {
	return b.Type
}

type SourceMessage struct {
	BaseImp
	Time uint64 `json:"time,omitempty"`
	ID   uint64 `json:"id,omitempty"`
}

type PlainMessage struct {
	BaseImp
	Text string `json:"text,omitempty"`
}

type AtMessage struct {
	BaseImp
	Target  uint64 `json:"target"`
	Display string `json:"display"`
}

type QuoteMessage struct {
	BaseImp
	Id       uint64 `json:"id"`
	GroupId  uint64 `json:"groupId"`
	SenderId uint64 `json:"senderId"`
	TargetId uint64 `json:"targetId"`
	Origin   Chain  `json:"origin"`
}
