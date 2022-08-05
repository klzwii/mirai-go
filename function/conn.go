package function

import (
	jsoniter "github.com/json-iterator/go"
	"sync/atomic"
)

type Conn interface {
	SendRequest(command string, subCommand string, req interface{}) error
}

type wsConn interface {
	WriteJSON(v interface{}) error
}

func GetWsConn(conn wsConn) Conn {
	return &connWsImp{
		conn: conn,
	}
}

type connWsImp struct {
	conn   wsConn
	syncID atomic.Uint32
}

type wsRequest struct {
	SyncId     uint32      `json:"syncId"`
	Command    string      `json:"command"`
	SubCommand string      `json:"subCommand,omitempty"`
	Content    interface{} `json:"content"`
}

func (c *connWsImp) SendRequest(command string, subCommand string, req interface{}) error {
	t := &wsRequest{
		SyncId:     c.syncID.Add(1),
		Command:    command,
		SubCommand: subCommand,
		Content:    req,
	}
	println(jsoniter.MarshalToString(t))
	return c.conn.WriteJSON(t)
}
