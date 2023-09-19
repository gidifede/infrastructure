package eventstore

//go:generate mockgen -source=./IEventStore.go -destination=./IEventStoreMock.go -package=eventStore

import "context"

type IEventStore interface {
	Append(context.Context, Event) error
	Get(context.Context, string) ([]Event, error)
}
