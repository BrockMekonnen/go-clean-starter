package events

type EventAddress struct {
	EventType string
	Topic     string
}

type Event[P any] struct {
	EventAddress
	EventID string
	Payload P
}

func (e Event[P]) GetEventAddress() EventAddress {
	return e.EventAddress
}

type EventInterface interface {
	GetEventAddress() EventAddress
}