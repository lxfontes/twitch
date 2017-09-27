package plugins

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/lxfontes/twitch/pkg/bot"
	"github.com/lxfontes/twitch/pkg/public"
)

const (
	songCommand = "!song"
)

type spotify struct {
	NoContext
	bot bot.Client
}

var _ Plugin = &spotify{}
var _ bot.MessageHandler = &spotify{}

func (g *spotify) Configure(b bot.Client, s *public.Server) {
	b.AddMessageHandler(songCommand, g)
	g.bot = b
}

func (g *spotify) OnMessage(from string, msg string) {
	// mac os specific .....
	var err error
	artist := ""
	album := ""
	track := ""

	if artist, err = osaScript(`tell application "Spotify" to artist of current track as string`); err != nil {
		log.Println(err)
		return
	}

	if album, err = osaScript(`tell application "Spotify" to album of current track as string`); err != nil {
		log.Println(err)
		return
	}

	if track, err = osaScript(`tell application "Spotify" to name of current track as string`); err != nil {
		log.Println(err)
		return
	}

	outMsg := fmt.Sprintf("Playing %s [%s] %s", artist, album, track)
	g.bot.Say(outMsg)
	fmt.Println(outMsg)
}

func osaScript(script string) (string, error) {
	args := []string{
		"-e",
		script,
	}

	cmd := exec.Command("osascript", args...)
	rawOut, err := cmd.Output()
	return strings.Trim(string(rawOut), "\n"), err
}
