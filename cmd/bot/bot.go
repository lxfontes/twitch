package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lxfontes/twitch/pkg/bot"
	"github.com/lxfontes/twitch/pkg/plugins"
	"github.com/lxfontes/twitch/pkg/public"
)

func main() {
	botUser := os.Getenv("BOT_USER")
	botPass := os.Getenv("BOT_PASS")
	botChannel := os.Getenv("BOT_CHANNEL")
	c, err := bot.NewClient(botUser, botPass, botChannel)
	if err != nil {
		panic(err)
	}

	s, err := public.NewServer(":8080")
	if err != nil {
		panic(err)
	}

	pl := plugins.DefaultStack()
	for _, plugin := range pl {
		plugin.Configure(c, s)
	}

	ctx, cancel := context.WithCancel(context.Background())
	for _, plugin := range pl {
		go plugin.Run(ctx)
	}

	go s.Serve()
	go c.Serve()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs

	cancel()
	c.Stop()
	s.Stop()
	log.Println("shutdown complete")
}
