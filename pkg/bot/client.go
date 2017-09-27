package bot

import (
	"crypto/tls"
	"log"
	"strings"
	"time"

	"github.com/lxfontes/twitch/pkg/store"
	irc "github.com/thoj/go-ircevent"
)

const (
	ircAddress = "irc.twitch.tv:6667"
)

type Client interface {
	AddMessageHandler(command string, handler MessageHandler)
	AddPresenceHandler(handler PresenceHandler)
	Say(msg string)
	Serve() error
	Stop()
}

type client struct {
	client           *irc.Connection
	channelName      string
	handlers         map[string]MessageHandler
	presenceHandlers []PresenceHandler
	store            store.Store
}

var _ Client = &client{}

type PresenceHandler interface {
	OnJoined(nick string)
	OnLeft(nick string)
}

type MessageHandler interface {
	OnMessage(from, msg string)
}

type HandlerFunc func(from, msg string)

func (hf HandlerFunc) Run(from, msg string) {
	hf(from, msg)
}

var _ Client = &client{}

func NewClient(nick, password, channel string, s store.Store) (*client, error) {
	ret := &client{
		channelName: formatChannel(channel),
		handlers:    map[string]MessageHandler{},
		store:       s,
	}

	c := irc.IRC(nick, "lxfB0T")
	c.Password = password
	c.VerboseCallbackHandler = true
	c.Debug = false
	c.UseTLS = false
	c.TLSConfig = &tls.Config{}
	c.AddCallback("001", func(e *irc.Event) {
		c.SendRaw("CAP REQ :twitch.tv/membership")
		c.Join(ret.channelName)
	})
	c.AddCallback("PRIVMSG", ret.handleMessage)
	c.AddCallback("JOIN", ret.handleJoin)
	c.AddCallback("PART", ret.handlePart)

	ret.client = c

	return ret, nil
}

func (c *client) handleJoin(e *irc.Event) {
	roster := c.store.Namespace("roster")
	u := User{}
	now := time.Now().Unix()

	err := roster.FindByID(e.Nick, &u)
	if err != nil {
		if err != store.ErrItemNotFound {
			log.Println(err)
			return
		}

		log.Println("first time seeing", e.Nick)

		u.CreatedAt = now
		u.Nick = e.Nick
	}

	u.UpdatedAt = now

	err = roster.Save(e.Nick, &u, store.NeverExpire)
	if err != nil {
		log.Println(err)
		return
	}

	// make sure we save our internal state first
	// as handlers might query the user data
	for _, handler := range c.presenceHandlers {
		handler.OnJoined(e.Nick)
	}
}

func (c *client) handlePart(e *irc.Event) {
	roster := c.store.Namespace("roster")
	u := User{}
	now := time.Now().Unix()

	err := roster.FindByID(e.Nick, &u)
	if err != nil {
		log.Println(err)
	}

	u.UpdatedAt = now

	err = roster.Save(e.Nick, &u, store.NeverExpire)
	if err != nil {
		log.Println(err)
		return
	}

	for _, handler := range c.presenceHandlers {
		handler.OnJoined(e.Nick)
	}
}

func (c *client) handleMessage(e *irc.Event) {
	user := e.User
	msg := e.Message()

	for prefix, h := range c.handlers {
		if strings.HasPrefix(msg, prefix) {
			h.OnMessage(user, msg)
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

func (c *client) AddMessageHandler(command string, handler MessageHandler) {
	c.handlers[command] = handler
}

func (c *client) AddPresenceHandler(handler PresenceHandler) {
	c.presenceHandlers = append(c.presenceHandlers, handler)
}

func formatChannel(s string) string {
	return "#" + s
}
