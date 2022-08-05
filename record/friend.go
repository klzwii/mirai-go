package record

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/entity"
)

type FriendMessageData struct {
	BaseDataImp
	friendMessageDataInner
}

type friendMessageDataInner struct {
	Sender entity.IndividualSender `json:"sender,omitempty"`
}

func GetFriendMessageRecord() Base {
	return &BaseImp{Data: &FriendMessageData{BaseDataImp: BaseDataImp{Type: FriendMessage}}}
}

func (f *FriendMessageData) UnmarshalJSON(data []byte) error {
	if err := jsoniter.Unmarshal(data, &f.BaseDataImp); err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, &f.friendMessageDataInner)
}
