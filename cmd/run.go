package cmd

import (
	"context"
	"fmt"
	bot2 "github.com/klzwii/mirai-go/bot"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

var isDebug = false

var rootCmd = &cobra.Command{
	Use:   "mirai-go",
	Short: "mirai-go is a go version mirai api frame work",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)
		bot, err := bot2.GetBot()
		if err != nil {
			log.Fatal(err)
		}
		bot.RegisterPlugin(&bot2.TestPlugin{})
		ctx, can := context.WithCancel(context.Background())
		defer can()
		if isDebug {
			log.SetLevel(log.DebugLevel)
			file, err := os.Create("./cpu.pprof")
			if err != nil {
				fmt.Printf("create cpu pprof failed, err:%v\n", err)
				return
			}
			_ = pprof.StartCPUProfile(file)
			defer pprof.StopCPUProfile()
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGUSR1)
			go func() {
				select {
				case _ = <-ch:
					can()
				}
			}()
		}
		bot.Start(ctx)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "start debugging")
}
