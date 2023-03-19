package message

import (
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"testing"
)

func TestGetMessagePlain(t *testing.T) {
	rawJson := "{\n                \"type\": \"Plain\",\n                \"text\": \"这个小酥肉好吃吗？\"\n            }"
	result := gjson.Parse(rawJson)
	message, err := getMessage(result)
	assert.Nil(t, err)
	plainMessage, ok := message.(*PlainMessage)
	assert.True(t, ok)
	assert.Equal(t, plainMessage.Type, PLAIN)
	assert.Equal(t, plainMessage.Text, "这个小酥肉好吃吗？")
}

func TestGetMessageUnknown(t *testing.T) {
	rawJson := "{\n                \"type\": \"Unknown\",\n                \"text\": \"这个小酥肉好吃吗？\"\n            }"
	result := gjson.Parse(rawJson)
	message, err := getMessage(result)
	assert.Nil(t, err)
	baseMessage, ok := message.(*BaseImp)
	assert.True(t, ok)
	assert.Equal(t, baseMessage.Type, Type("Unknown"))
}

func TestGetMessageSource(t *testing.T) {
	rawJson := "{\n                \"type\": \"Source\",\n                \"id\": 1422833,\n                \"time\": 1659604654\n            }"
	result := gjson.Parse(rawJson)
	message, err := getMessage(result)
	assert.Nil(t, err)
	sourceMessage, ok := message.(*SourceMessage)
	assert.True(t, ok)
	assert.Equal(t, sourceMessage.Type, SOURCE)
	assert.Equal(t, sourceMessage.Time, uint64(1659604654))
	assert.Equal(t, sourceMessage.ID, uint64(1422833))
}
