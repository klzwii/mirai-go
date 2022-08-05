package main

import (
	"github.com/gorilla/websocket"
	"github.com/klzwii/mirai-go/assembler"
	"log"
	"net/url"
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
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		assembler.MarshalToRecord(string(message))
		log.Printf("recv: %s", message)
	}
}
