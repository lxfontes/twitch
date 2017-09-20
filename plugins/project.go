package plugins

import (
	"net/http"

	"github.com/lxfontes/twitch/pkg/bot"
	"github.com/lxfontes/twitch/public"
)

type project struct {
	NoContext
	bot bot.Client
}

var _ Plugin = &project{}
var _ bot.Handler = &project{}
var _ http.Handler = &project{}

func (p *project) Configure(b bot.Client, s *public.Server) {
	b.Register("!project", p)
	s.Router().Handle("/project", p)
	p.bot = b
}

func (p *project) ServeIRC(from string, msg string) {
	p.bot.Say("luls")
}

func (p *project) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("damn"))
}
