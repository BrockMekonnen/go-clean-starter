package events

func MakeEventConsumer(subscriber Subscriber, address EventAddress, fn func() func(EventInterface) error, opts *SubscriberOptions) func() error {
	return func() error {
		handler := fn()
		return subscriber.Add(address, handler, opts)
	}
}
