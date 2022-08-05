package sender

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockConn struct {
}

func (m *MockConn) SendRequest(command string, subCommand string, req interface{}) error {
	println(jsoniter.MarshalToString(req))
	return nil
}

func TestSenderImp_SendToFriend(t *testing.T) {
	sender := GetWSSender(&MockConn{}, "")
	assert.Nil(t, sender.SendToFriend())
}
