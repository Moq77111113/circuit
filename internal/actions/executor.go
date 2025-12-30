package actions

import (
	"context"
	"time"
)

const DefaultTimeout = 30 * time.Second

type Def struct {
	Name                string
	Label               string
	Description         string
	Run                 func(context.Context) error
	Timeout             time.Duration
	RequireConfirmation bool
}

func Execute(ctx context.Context, action Def) error {
	timeout := action.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return action.Run(ctx)
}
