package function

import (
	"github.com/klzwii/mirai-go/record"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockConn struct {
	testFunc func(v interface{})
}

func (m *mockConn) ReadMessage() (messageType int, p []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockConn) WriteJSON(v interface{}) error {
	m.testFunc(v)
	return nil
}

var testData = map[string]string{
	"asd": "123",
	"asc": "346",
}

func TestConnWsImp_SendRequest(t *testing.T) {
	conn := &mockConn{}
	wsCon := GetWsConn(conn)
	conn.testFunc = func(v interface{}) {
		req, ok := v.(*wsRequest)
		assert.True(t, ok)
		assert.Equal(t, uint32(1), req.SyncId)
		assert.Equal(t, "testMainCommand", req.Command)
		assert.Equal(t, "testSubCommand", req.SubCommand)
		assert.Equal(t, testData, req.Content)
		_ = wsCon.(*connWsImp).center.Notify(req.SyncId, []byte("{}"), nil)
	}
	err := wsCon.SendRequest("testMainCommand", "testSubCommand", testData, record.GetSendMessageResponse())
	assert.Nil(t, err)
}
