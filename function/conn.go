package function

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/klzwii/mirai-go/assembler"
	"github.com/klzwii/mirai-go/record"
	"github.com/klzwii/mirai-go/util"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
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
	SendRequest(command string, subCommand string, req any, resp any) error
	StartReading(ctx context.Context, ch chan record.Base)
}

type wsConn interface {
	WriteJSON(v interface{}) error
	ReadMessage() (messageType int, p []byte, err error)
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

func (c *connWsImp) SendRequest(command string, subCommand string, req any, resp any) error {
	syncId, ch := c.center.RegisterEvent()
	t := &wsRequest{
		SyncId:     syncId,
		Command:    command,
		SubCommand: subCommand,
		Content:    req,
	}
	log.Debugf("send ws request %+v", *t)
	if err := c.conn.WriteJSON(t); err != nil {
		return err
	}
	timer := timerPool.Get().(*time.Timer)
	timer.Reset(20 * time.Second)
	defer func() {
		timer.Stop()
		timerPool.Put(timer)
	}()
	start := time.Now().UnixMicro()
	select {
	case ret := <-ch:
		log.Debug("get response after ", time.Now().UnixMicro()-start)
		if ret.Err != nil {
			return ret.Err
		}

		return json.Unmarshal(ret.Data.([]byte), resp)
	case <-timer.C:
		return ErrTimeOut
	}
}

func (c *connWsImp) StartReading(ctx context.Context, ch chan record.Base) {
	callback := getDispatchFunc(c.center, ch)
	for {
		select {
		case <-ctx.Done():
			break
		default:
			_, m, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Debugf("Get ws message %v", string(m))
			go callback(m)
		}
	}
}

func getDispatchFunc(e util.EventCenter, ch chan record.Base) func(message []byte) {
	return func(message []byte) {
		syncID := gjson.GetBytes(message, "syncId").Str
		log.Debug("sync ID is ", syncID)
		if syncID == "-1" {
			ch <- assembler.UnmarshalToRecord(message)
		} else if len(syncID) != 0 {
			var (
				eventID uint64
				err     error
			)
			if eventID, err = strconv.ParseUint(syncID, 10, 64); err != nil {
				log.Error("parse sync ID failed ", err)
				return
			}
			_ = e.Notify(uint32(eventID), message)
			log.Debug("notified event ", eventID)
		}
	}
}
