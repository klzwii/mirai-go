package record

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/entity"
)

type FriendMessageData struct {
	BaseData
	friendMessageDataInner
}

type friendMessageDataInner struct {
	Sender entity.IndividualSender `json:"sender,omitempty"`
}

func GetFriendMessageRecord() BaseInterface {
	return &Base{Data: &FriendMessageData{BaseData: BaseData{Type: FriendMessage}}}
}

func (f *FriendMessageData) UnmarshalJSON(data []byte) error {
	if err := jsoniter.Unmarshal(data, &f.BaseData); err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, &f.friendMessageDataInner)
}
