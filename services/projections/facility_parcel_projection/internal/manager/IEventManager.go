package manager

import "context"

type IEventManager interface {
	ManageEvent(ctx context.Context) error
}
