package command

import "context"

// Command ...
type Command interface {
	Run(ctx context.Context, args []string) error
}
