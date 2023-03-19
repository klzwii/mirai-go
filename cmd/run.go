package cmd

import (
	"context"
	"fmt"
	bot2 "github.com/klzwii/mirai-go/bot"
	"github.com/klzwii/mirai-go/chatgpt"
	"github.com/klzwii/mirai-go/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func loadConfig() (map[string]any, error) {
	file, err := os.Open(util.ConfigFile)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var ret map[string]any
	return ret, yaml.Unmarshal(data, &ret)
}

var rootCmd = &cobra.Command{
	Use:   "mirai-go",
	Short: "mirai-go is a go version mirai api frame work",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)
		var (
			config map[string]interface{}
			err    error
		)
		if config, err = loadConfig(); err != nil {
			log.Fatal(err)
		}
		var bot *bot2.Bot
		if bot, err = bot2.GetBot(config); err != nil {
			log.Fatal(err)
		}
		bot.RegisterPlugin(&bot2.TestPlugin{})
		bot.RegisterPlugin(&chatgpt.Plugin{})
		ctx, can := context.WithCancel(context.Background())
		defer can()
		if util.IsDebug {
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
	rootCmd.PersistentFlags().BoolVarP(&util.IsDebug, "debug", "d", false, "start debugging")
	rootCmd.PersistentFlags().StringVarP(&util.ConfigFile, "config", "c", "mirai_config.yaml", "set config file")
}
