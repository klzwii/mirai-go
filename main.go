package main

import (
	"context"
	bot2 "github.com/klzwii/mirai-go/bot"
)

func main() {
	bot, err := bot2.GetBot()
	if err != nil {
		panic(err)
	}
	bot.RegisterPlugin(&bot2.TestPlugin{})
	bot.Start(context.TODO())
}
