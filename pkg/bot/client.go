package bot

import (
	"crypto/tls"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

const (
	ircAddress = "irc.twitch.tv:6667"
)

type Client interface {
	Register(command string, handler Handler)
	Say(msg string)
	Serve() error
	Stop()
}

type client struct {
	client      *irc.Connection
	channelName string
	handlers    map[string]Handler
}

var _ Client = &client{}

type Handler interface {
	ServeIRC(from, msg string)
}

type HandlerFunc func(from, msg string)

func (hf HandlerFunc) Run(from, msg string) {
	hf(from, msg)
}

var _ Client = &client{}

func NewClient(nick, password, channel string) (*client, error) {
	ret := &client{
		channelName: formatChannel(channel),
		handlers:    map[string]Handler{},
	}

	c := irc.IRC(nick, "lxfB0T")
	c.Password = password
	c.VerboseCallbackHandler = true
	c.Debug = true
	c.UseTLS = false
	c.TLSConfig = &tls.Config{}
	c.AddCallback("001", func(e *irc.Event) { c.Join(ret.channelName) })
	c.AddCallback("PRIVMSG", ret.handleMessage)

	ret.client = c

	return ret, nil
}

func (c *client) handleMessage(e *irc.Event) {
	user := e.User
	msg := e.Message()

	for prefix, h := range c.handlers {
		if strings.HasPrefix(msg, prefix) {
			h.ServeIRC(user, msg)
		}
	}
}

func (c *client) Serve() error {
	if err := c.client.Connect(ircAddress); err != nil {
		return err
	}

	c.client.Loop()
	return <-c.client.Error
}

func (c *client) Stop() {
	c.client.Quit()
}

func (c *client) Say(msg string) {
	c.client.Privmsg(c.channelName, msg)
}

func (c *client) Register(command string, handler Handler) {
	c.handlers[command] = handler
}

func formatChannel(s string) string {
	return "#" + s
}
