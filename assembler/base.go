package assembler

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/record"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type recordFunc func() record.Base

var (
	funcMap = map[record.Type]recordFunc{
		record.GroupMessage:  record.GetGroupMessageRecord,
		record.FriendMessage: record.GetFriendMessageRecord,
	}
	nilObject = &record.BaseImp{
		SyncID: "-2",
		Data:   &record.BaseDataImp{Type: record.NULL},
	}
)

func innerMarshal(rawJson []byte, ret record.Base) (record.Base, error) {
	if err := jsoniter.Unmarshal(rawJson, ret); err != nil {
		log.Println(err)
		return nil, err
	}
	return ret, nil
}

// UnmarshalToRecord convert raw json Message to a record
func UnmarshalToRecord(rawMessage []byte) record.Base {
	var recordType = record.NULL
	data := gjson.GetBytes(rawMessage, "data")
	if data.IsObject() {
		recordType = record.ConvertToType(data.Get("type").Str)
	}
	if fun, ok := funcMap[recordType]; ok {
		if ret, err := innerMarshal(rawMessage, fun()); err == nil {
			return ret
		}
	}
	return nilObject
}
