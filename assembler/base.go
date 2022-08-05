package assembler

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/record"
	"github.com/tidwall/gjson"
	"log"
)

type recordFunc func() record.BaseInterface

var (
	logger  = log.Default()
	funcMap = map[record.Type]recordFunc{
		record.GroupMessage:  record.GetGroupMessageRecord,
		record.FriendMessage: record.GetFriendMessageRecord,
	}
	nilObject = &record.Base{
		SyncID: "-2",
		Data:   &record.BaseData{Type: record.NULL},
	}
)

func innerMarshal(rawJson string, ret record.BaseInterface) (record.BaseInterface, error) {
	if err := jsoniter.UnmarshalFromString(rawJson, ret); err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// MarshalToRecord convert raw json Message
func MarshalToRecord(rawMessage string) record.BaseInterface {
	syncID := gjson.Get(rawMessage, "syncID").Str
	logger.Println("sync_id is " + syncID)
	var recordType record.Type = record.NULL
	data := gjson.Get(rawMessage, "data")
	if data.IsObject() {
		recordType = record.ConvertToType(data.Get("type").Str)
	}
	logger.Println("record type is" + recordType)
	if fun, ok := funcMap[recordType]; ok {
		if ret, err := innerMarshal(rawMessage, fun()); err != nil {
			return ret
		}
	}
	return nilObject
}
