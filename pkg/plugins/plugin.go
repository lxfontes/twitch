package plugins

import (
	"context"

	"github.com/lxfontes/twitch/pkg/bot"
	"github.com/lxfontes/twitch/pkg/public"
)

func DefaultStack() []Plugin {
	return []Plugin{
		&giffer{},
	}
}

type Plugin interface {
	Configure(bot.Client, *public.Server)
	Run(ctx context.Context)
}

type NoContext struct{}

func (n *NoContext) Run(_ context.Context) {}
