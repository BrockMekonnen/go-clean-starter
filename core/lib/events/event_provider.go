package events

import (
	"context"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/contracts"
)

// EnqueueFunc is like the TS Enqueue<E extends Event>
type EnqueueFunc func(EventInterface)

// EventStore holds and manages queued events
type eventStore struct {
	events []EventInterface
}

func newEventStore() *eventStore {
	return &eventStore{events: make([]EventInterface, 0)}
}

func (es *eventStore) Enqueue(event EventInterface) {
	es.events = append(es.events, event)
}

func (es *eventStore) GetEvents() []EventInterface {
	return es.events
}

// EventUsecaseFactory mirrors: (deps, enqueue) => ApplicationService
type EventUsecaseFactory[Deps any, Payload any, Result any] func(deps Deps, enqueue EnqueueFunc) contracts.ApplicationService[Payload, Result]

// MakeEventProvider wraps a usecase factory to defer event publishing
func EventProvider[Deps any, Payload any, Result any](
	factory EventUsecaseFactory[Deps, Payload, Result],
) func(deps Deps, publisher Publisher) contracts.ApplicationService[Payload, Result] {
	return func(deps Deps, publisher Publisher) contracts.ApplicationService[Payload, Result] {
		es := newEventStore()
		usecase := factory(deps, es.Enqueue)

		return func(ctx context.Context, payload Payload) (Result, error) {
			result, err := usecase(ctx, payload)
			if err != nil {
				return result, err
			}

			for _, event := range es.GetEvents() {
				publisher.Publish(event)
			}

			return result, nil
		}
	}
}
