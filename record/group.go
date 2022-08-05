package record

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/entity"
)

type GroupMessageData struct {
	BaseDataImp
	groupMessageDataInner
}

type groupMessageDataInner struct {
	Sender entity.GroupSender `json:"sender,omitempty"`
}

func GetGroupMessageRecord() Base {
	return &BaseImp{Data: &GroupMessageData{BaseDataImp: BaseDataImp{Type: GroupMessage}}}
}

func (g *GroupMessageData) UnmarshalJSON(data []byte) error {
	if err := jsoniter.Unmarshal(data, &g.BaseDataImp); err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, &g.groupMessageDataInner)
}
