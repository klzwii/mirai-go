package message

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

type Type string

const (
	SOURCE = Type("Source")
	PLAIN  = Type("Plain")
)

func GetMessage(result gjson.Result) (Base, error) {
	messageType := Type(result.Get("type").Str)
	var ret Base = nil
	switch messageType {
	case SOURCE:
		ret = &SourceMessage{}
	case PLAIN:
		ret = &PlainMessage{}
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
