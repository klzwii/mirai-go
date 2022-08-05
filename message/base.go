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

func GetMessage(result gjson.Result) (BaseInterface, error) {
	messageType := Type(result.Get("type").Str)
	var ret BaseInterface = nil
	switch messageType {
	case SOURCE:
		ret = &SourceMessage{}
	case PLAIN:
		ret = &PlainMessage{}
	default:
		ret = &Base{}
	}
	if err := jsoniter.UnmarshalFromString(result.Raw, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type BaseInterface interface {
	GetType() Type
}

type Base struct {
	Type Type `json:"type,omitempty"`
}

func (b *Base) GetType() Type {
	return b.Type
}

type SourceMessage struct {
	Base
	Time uint64 `json:"time,omitempty"`
	ID   uint64 `json:"id,omitempty"`
}

type PlainMessage struct {
	Base
	Text string `json:"text,omitempty"`
}
