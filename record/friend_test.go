package record

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/klzwii/mirai-go/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFriendMessageData(t *testing.T) {
	record := GetFriendMessageRecord()
	err := jsoniter.UnmarshalFromString("{\n    \"syncId\": \"-1\",\n    \"data\": {\n        \"type\": \"FriendMessage\",\n        \"messageChain\": [\n            {\n                \"type\": \"Source\",\n                \"id\": 17705,\n                \"time\": 1659604789\n            },\n            {\n                \"type\": \"Plain\",\n                \"text\": \"hi\"\n            }\n        ],\n        \"sender\": {\n            \"id\": 1027898733,\n            \"nickname\": \"klzwii(已黑化)\",\n            \"remark\": \"klzwii\"\n        }\n    }\n}", record)
	assert.Nil(t, err)
	assert.Equal(t, record.GetSyncID(), "-1")

	friendMessageData, ok := record.GetData().(*FriendMessageData)
	assert.True(t, ok)

	chain := friendMessageData.MessageChain
	assert.Equal(t, 2, len(chain))

	assert.Equal(t, message.SOURCE, chain[0].GetType())
	sourceMessage, ok := chain[0].(*message.SourceMessage)
	assert.True(t, ok)
	assert.Equal(t, uint64(17705), sourceMessage.ID)
	assert.Equal(t, uint64(1659604789), sourceMessage.Time)

	assert.Equal(t, message.PLAIN, chain[1].GetType())
	plainMessage, ok := chain[1].(*message.PlainMessage)
	assert.True(t, ok)
	assert.Equal(t, "hi", plainMessage.Text)

	sender := friendMessageData.Sender
	assert.Equal(t, uint64(1027898733), sender.ID)
	assert.Equal(t, "klzwii(已黑化)", sender.Nickname)
	assert.Equal(t, "klzwii", sender.Remark)
}
