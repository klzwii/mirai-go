package function

import (
	"errors"
	"github.com/klzwii/mirai-go/util"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var timerPool = sync.Pool{
	New: func() any {
		timer := time.NewTimer(10 * time.Second)
		timer.Stop()
		return timer
	},
}

type Conn interface {
	SendRequest(command string, subCommand string, req interface{}) ([]byte, error)
}

type wsConn interface {
	WriteJSON(v interface{}) error
}

func GetWsConn(conn wsConn) Conn {
	return &connWsImp{
		conn:   conn,
		center: util.New(2000),
	}
}

type connWsImp struct {
	conn   wsConn
	center util.EventCenter
}

type wsRequest struct {
	SyncId     uint32      `json:"syncId"`
	Command    string      `json:"command"`
	SubCommand string      `json:"subCommand,omitempty"`
	Content    interface{} `json:"content"`
}

var ErrTimeOut = errors.New("websocket send request timeout")

func (c *connWsImp) SendRequest(command string, subCommand string, req interface{}) ([]byte, error) {
	syncId, ch := c.center.RegisterEvent()
	t := &wsRequest{
		SyncId:     syncId,
		Command:    command,
		SubCommand: subCommand,
		Content:    req,
	}
	log.Debugf("send ws request %+v", *t)
	if err := c.conn.WriteJSON(t); err != nil {
		return nil, err
	}
	timer := timerPool.Get().(*time.Timer)
	timer.Reset(20 * time.Second)
	defer timerPool.Put(timer)
	select {
	case ret := <-ch:
		return ret.Data.([]byte), ret.Err
	case <-timer.C:
		return nil, ErrTimeOut
	}
}
