package bot

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/gorilla/websocket"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/sender"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetBot(t *testing.T) {
	patch := gomonkey.NewPatches()
	defer patch.Reset()
	testConfig := map[string]any{
		"bot": map[string]any{
			"qq":         123,
			"port":       456,
			"verify_key": 789,
		},
	}
	var mockDialer = &websocket.Dialer{}
	patch.ApplyMethod(mockDialer, "Dial", func(_ *websocket.Dialer, s string,
		data http.Header) (*websocket.Conn, *http.Response, error) {
		assert.Equal(t, "ws://localhost:456/message", s)
		assert.Equal(t, "123", data["qq"][0])
		assert.Equal(t, "789", data["verifyKey"][0])
		return nil, nil, nil
	})
	_, err := GetBot(testConfig)
	if err != nil {
		panic(err)
	}
}

type MockPlugin struct {
	t *testing.T
}

func (m *MockPlugin) OnGroupMessage(_ *record.GroupMessageData) {
	panic("implement me")
}

func (m *MockPlugin) OnFriendMessage(_ *record.FriendMessageData) {
	panic("implement me")
}

func (m *MockPlugin) Name() string {
	return "TestPlugin"
}

func (m *MockPlugin) Initiate(conf map[string]any, _ sender.Sender) {
	assert.Equal(m.t, "test", conf["data"])
}

func Test_initiatePlugins(t *testing.T) {
	initiatePlugins(map[string]any{
		"TestPlugin": map[string]any{
			"data": "test",
		},
	}, []Plugin{&MockPlugin{t: t}}, nil)
}
