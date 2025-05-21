package events

func MakeEventConsumer[D any](subscriber Subscriber, address EventAddress, fn func(D) func(EventInterface) error, opts *SubscriberOptions) func(D) error {
	return func(deps D) error {
		handler := fn(deps)
		return subscriber.Add(address, handler, opts)
	}
}