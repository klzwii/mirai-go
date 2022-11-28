package assembler

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/record"
	"github.com/tidwall/gjson"
	"log"
)

type recordFunc func() record.Base

var (
	logger  = log.Default()
	funcMap = map[record.Type]recordFunc{
		record.GroupMessage:  record.GetGroupMessageRecord,
		record.FriendMessage: record.GetFriendMessageRecord,
	}
	nilObject = &record.BaseImp{
		SyncID: "-2",
		Data:   &record.BaseDataImp{Type: record.NULL},
	}
)

func innerMarshal(rawJson string, ret record.Base) (record.Base, error) {
	if err := jsoniter.UnmarshalFromString(rawJson, ret); err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// UnmarshalToRecord convert raw json Message to a record
func UnmarshalToRecord(rawMessage string) record.Base {
	//syncID := gjson.Get(rawMessage, "syncID").Str
	//logger.Println("sync_id is " + syncID)
	var recordType = record.NULL
	data := gjson.Get(rawMessage, "data")
	if data.IsObject() {
		recordType = record.ConvertToType(data.Get("type").Str)
	}
	logger.Println("record type is" + recordType)
	if fun, ok := funcMap[recordType]; ok {
		if ret, err := innerMarshal(rawMessage, fun()); err == nil {
			return ret
		}
	}
	return nilObject
}
