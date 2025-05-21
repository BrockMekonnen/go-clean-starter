// Updated pubsub.go with channel-based async dispatcher
package events

import (
	"fmt"
	"sync"
)

type SubscriberOptions struct {
	Single bool
	NackOn func(error) bool
}

type Publisher interface {
	Publish(event EventInterface) error
}

type Subscriber interface {
	Add(address EventAddress, handler func(EventInterface) error, opts *SubscriberOptions) error
	Start() error
	Dispose() error
}

type EventEmitterPubSub struct {
	mu        sync.Mutex
	listeners map[string][]handlerRegistration
	started   bool
	defaults  SubscriberOptions

	eventChan chan EventInterface
	quitChan  chan struct{}
	wg        sync.WaitGroup
}

type handlerRegistration struct {
	single  bool
	handler func(EventInterface) error
}

func NewEventEmitterPubSub() *EventEmitterPubSub {
	e := &EventEmitterPubSub{
		listeners: make(map[string][]handlerRegistration),
		defaults: SubscriberOptions{
			Single: false,
			NackOn: func(err error) bool { return true },
		},
		eventChan: make(chan EventInterface, 100), // buffer size
		quitChan:  make(chan struct{}),
	}
	return e
}

func (e *EventEmitterPubSub) key(addr EventAddress) string {
	return fmt.Sprintf("%s.%s", addr.Topic, addr.EventType)
}

func (e *EventEmitterPubSub) Publish(event EventInterface) error {
	select {
	case e.eventChan <- event:
		return nil
	default:
		return fmt.Errorf("event channel is full, dropping event")
	}
}

func (e *EventEmitterPubSub) Add(address EventAddress, handler func(EventInterface) error, opts *SubscriberOptions) error {
	if opts == nil {
		opts = &e.defaults
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	key := e.key(address)
	e.listeners[key] = append(e.listeners[key], handlerRegistration{
		single:  opts.Single,
		handler: handler,
	})

	return nil
}

func (e *EventEmitterPubSub) Start() error {
	if e.started {
		return nil
	}
	e.started = true
	e.wg.Add(1)
	go e.dispatchLoop()
	return nil
}

func (e *EventEmitterPubSub) dispatchLoop() {
	defer e.wg.Done()
	for {
		select {
		case event := <-e.eventChan:
			e.dispatch(event)
		case <-e.quitChan:
			return
		}
	}
}

func (e *EventEmitterPubSub) dispatch(event EventInterface) {
	key := e.key(event.GetEventAddress())
	e.mu.Lock()
	handlers := e.listeners[key]

	var remaining []handlerRegistration
	for _, reg := range handlers {
		if !reg.single {
			remaining = append(remaining, reg)
		}
	}
	if len(remaining) > 0 {
		e.listeners[key] = remaining
	} else {
		delete(e.listeners, key)
	}
	e.mu.Unlock()

	for _, reg := range handlers {
		go func(r handlerRegistration) {
			if err := r.handler(event); err != nil && r.single && e.defaults.NackOn(err) {
				fmt.Printf("event handler error: %v\n", err)
			}
		}(reg)
	}
}

func (e *EventEmitterPubSub) Dispose() error {
	if !e.started {
		return nil
	}
	close(e.quitChan)
	e.wg.Wait()

	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners = make(map[string][]handlerRegistration)
	e.started = false
	return nil
}