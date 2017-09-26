package plugins

import (
	"log"
	"net/http"
	"os"
	"strings"

	giphy "github.com/ivolo/go-giphy"
	"github.com/lxfontes/twitch/pkg/bot"
	"github.com/lxfontes/twitch/pkg/public"
)

const (
	gifCommand = "!gif"
)

type giffer struct {
	NoContext
	bot    bot.Client
	public *public.Server
	client *giphy.Client
}

var _ Plugin = &giffer{}
var _ bot.Handler = &giffer{}
var _ http.Handler = &giffer{}

func (g *giffer) Configure(b bot.Client, s *public.Server) {
	b.Register(gifCommand, g)
	s.OBSScript("/js/gif.js")
	s.Router().Handle("/event/gif", g)
	g.bot = b
	g.public = s
	g.client = giphy.New(os.Getenv("BOT_GIPHY_KEY"))
}

func (g *giffer) ServeIRC(from string, msg string) {
	what := strings.TrimSpace(strings.TrimPrefix(msg, gifCommand))
	if what == "" {
		log.Println("blank gif query")
		return
	}

	log.Println("gif: searching for", what)

	gifs, err := g.client.Search(what)
	if err != nil {
		log.Println("giffer:", err)
		return
	}

	if len(gifs) == 0 {
		log.Println("giffer: nothing found")
	}

	gif := public.NewCommand("gif")
	gif.Arguments["target"] = gifs[0].Images["original"].URL
	g.public.SendCommand(gif)
}

func (g *giffer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gif := public.NewCommand("gif")
	gif.Arguments["target"] = "https://media.giphy.com/media/3owvK6fi98W6lCKZSE/giphy.gif"
	g.public.SendCommand(gif)
}
