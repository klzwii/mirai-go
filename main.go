package main

import (
	"github.com/gorilla/websocket"
	"github.com/klzwii/mirai-go/assembler"
	"github.com/klzwii/mirai-go/function"
	"github.com/klzwii/mirai-go/function/sender"
	"github.com/klzwii/mirai-go/message"
	"log"
	"net/url"
	"time"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:5679", Path: "/message"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), map[string][]string{
		"verifyKey": {"1234567890"},
		"qq":        {"2844255154"},
	})
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		senderImp := sender.GetWSSender(function.GetWsConn(c), "")
		for i := 0; i < 4; i++ {
			time.Sleep(3 * time.Second)
			if err := senderImp.SendToGroup(590258464, message.NewMessageChain().AddPlain("123")); err != nil {
				log.Fatal(err)
			}
		}
	}()
	for {
		_, m, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		assembler.MarshalToRecord(string(m))
		log.Printf("recv: %s", m)
	}
}
