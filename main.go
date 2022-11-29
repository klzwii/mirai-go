package main

import (
	"context"
	"fmt"
	bot2 "github.com/klzwii/mirai-go/bot"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func main() {
	file, err := os.Create("./cpu.pprof")
	if err != nil {
		fmt.Printf("create cpu pprof failed, err:%v\n", err)
		return
	}

	log.SetLevel(log.DebugLevel)
	bot, err := bot2.GetBot()
	if err != nil {
		panic(err)
	}
	bot.RegisterPlugin(&bot2.TestPlugin{})
	ctx, can := context.WithCancel(context.Background())
	{
		_ = pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)
		select {
		case _ = <-ch:
			can()
		}
	}
	bot.Start(ctx)
}
