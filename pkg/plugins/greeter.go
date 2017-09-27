package plugins

import (
	"fmt"

	"github.com/lxfontes/twitch/pkg/bot"
	"github.com/lxfontes/twitch/pkg/public"
)

type greeter struct {
	NoContext
	bot bot.Client
}

var _ Plugin = &greeter{}
var _ bot.PresenceHandler = &greeter{}

func (g *greeter) Configure(b bot.Client, s *public.Server) {
	b.AddPresenceHandler(g)
	g.bot = b
}

func (g *greeter) OnJoined(nick string) {
	g.bot.Say(fmt.Sprintf("ohai %s !!", nick))
}

// noop
func (g *greeter) OnLeft(nick string) {}
