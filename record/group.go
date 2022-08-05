package record

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/entity"
)

type GroupMessageData struct {
	BaseData
	groupMessageDataInner
}

type groupMessageDataInner struct {
	Sender entity.GroupSender `json:"sender,omitempty"`
}

func GetGroupMessageRecord() BaseInterface {
	return &Base{Data: &GroupMessageData{BaseData: BaseData{Type: GroupMessage}}}
}

func (g *GroupMessageData) UnmarshalJSON(data []byte) error {
	if err := jsoniter.Unmarshal(data, &g.BaseData); err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, &g.groupMessageDataInner)
}
