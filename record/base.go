package record

import (
	"errors"
	"github.com/klzwii/mirai-go/message"
	"github.com/tidwall/gjson"
)

type Type string

const (
	NULL                = Type("")
	FriendMessage       = Type("FriendMessage")
	GroupMessage        = Type("GroupMessage")
	SendMessageResponse = Type("SendMessageResponse")
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
	results := gjson.GetManyBytes(data, "type", "messageChain")
	if len(results) < 1 {
		return NotEnoughFieldsError
	}
	b.Type = Type(results[0].Str)
	if len(results) < 2 || !results[1].IsArray() {
		return nil
	}
	chain, err := message.GetMessageChain(results[1].Array())
	if err != nil {
		return err
	}
	b.MessageChain = chain
	return nil
}

func (b *BaseDataImp) getType() Type {
	return b.Type
}

type Base interface {
	GetSyncID() string
	GetData() BaseData
	GetType() Type
}

type BaseImp struct {
	SyncID string   `json:"syncId,omitempty"`
	Data   BaseData `json:"data,omitempty"`
}

func (b *BaseImp) GetSyncID() string {
	return b.SyncID
}

func (b *BaseImp) GetData() BaseData {
	return b.Data
}

func (b *BaseImp) GetType() Type {
	return b.Data.getType()
}
