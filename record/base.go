package record

import (
	"errors"
	"github.com/klzwii/mirai-go/message"
	"github.com/tidwall/gjson"
)

type Type string

const (
	NULL          = Type("")
	FriendMessage = Type("FriendMessage")
	GroupMessage  = Type("GroupMessage")
)

var (
	string2Type = map[string]Type{
		"FriendMessage": FriendMessage,
		"GroupMessage":  GroupMessage,
	}
	NotEnoughFieldsError = errors.New("insufficient field")
)

func ConvertToType(rawType string) Type {
	t, ok := string2Type[rawType]
	if !ok {
		return NULL
	}
	return t
}

type BaseData interface {
	getType() Type
}

type BaseDataImp struct {
	Type         Type          `json:"type,omitempty"`
	MessageChain message.Chain `json:"messageChain,omitempty"`
}

func (b *BaseDataImp) UnmarshalJSON(data []byte) error {
	results := gjson.GetMany(string(data), "type", "messageChain")
	if len(results) < 1 {
		return NotEnoughFieldsError
	}
	b.Type = Type(results[0].Str)
	if len(results) < 2 || !results[1].IsArray() {
		return nil
	}
	b.MessageChain = message.Chain{}
	for _, result := range results[1].Array() {
		if ele, err := message.GetMessage(result); err != nil {
			return err
		} else {
			b.MessageChain = append(b.MessageChain, ele)
		}
	}
	return nil
}

func (b *BaseDataImp) getType() Type {
	return b.Type
}

type Base interface {
	getSyncID() string
	getData() BaseData
}

type BaseImp struct {
	SyncID string   `json:"syncId,omitempty"`
	Data   BaseData `json:"data,omitempty"`
}

func (b *BaseImp) getSyncID() string {
	return b.SyncID
}

func (b *BaseImp) getData() BaseData {
	return b.Data
}
