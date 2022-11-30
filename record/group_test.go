package record

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGroupMessageRecord(t *testing.T) {
	record := GetGroupMessageRecord()
	err := jsoniter.Unmarshal([]byte("{\n    \"syncId\": \"-1\",\n    \"data\": {\n        \"type\": \"GroupMessage\",\n        \"messageChain\": [\n            {\n                \"type\": \"Source\",\n                \"id\": 1422833,\n                \"time\": 1659604654\n            },\n            {\n                \"type\": \"Plain\",\n                \"text\": \"这个小酥肉好吃吗？\"\n            }\n        ],\n        \"sender\": {\n            \"id\": 553127268,\n            \"memberName\": \"毛\",\n            \"specialTitle\": \"\",\n            \"permission\": \"MEMBER\",\n            \"joinTimestamp\": 1655306659,\n            \"lastSpeakTimestamp\": 1659604654,\n            \"muteTimeRemaining\": 0,\n            \"group\": {\n                \"id\": 398294930,\n                \"name\": \"東方緑豆汤 ~ Dream of the Green Dragon\",\n                \"permission\": \"MEMBER\"\n            }\n        }\n    }\n}"), record)
	assert.Nil(t, err)
	groupRecord, ok := record.(*BaseImp)
	assert.True(t, ok)
	groupData, ok := groupRecord.Data.(*GroupMessageData)
	assert.True(t, ok)
	assert.Equal(t, groupRecord.SyncID, "-1")
	// assert message chain
	assert.Equal(t, len(groupData.MessageChain), 2)
	assert.Equal(t, groupData.MessageChain[0].GetType(), message.SOURCE)
	sourceMessage, ok := groupData.MessageChain[0].(*message.SourceMessage)
	assert.True(t, ok)
	assert.Equal(t, sourceMessage.Time, uint64(1659604654))
	assert.Equal(t, sourceMessage.ID, uint64(1422833))
	assert.Equal(t, groupData.MessageChain[1].GetType(), message.PLAIN)
	plainMessage, ok := groupData.MessageChain[1].(*message.PlainMessage)
	assert.True(t, ok)
	assert.Equal(t, plainMessage.Text, "这个小酥肉好吃吗？")
	// assert sender
	assert.Equal(t, groupData.Sender.ID, uint64(553127268))
	assert.Equal(t, groupData.Sender.MemberName, "毛")
	assert.Equal(t, groupData.Sender.SpecialTitle, "")
	assert.Equal(t, groupData.Sender.Permission, "MEMBER")
	assert.Equal(t, groupData.Sender.JoinTimestamp, uint64(1655306659))
	assert.Equal(t, groupData.Sender.LastSpeakTimestamp, uint64(1659604654))
	assert.Equal(t, groupData.Sender.MuteTimeRemaining, uint64(0))
	assert.Equal(t, groupData.Sender.Group.ID, uint64(398294930))
}
