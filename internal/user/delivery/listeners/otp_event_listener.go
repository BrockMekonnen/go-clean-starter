package listeners

import (
	"fmt"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/events"
	app "github.com/BrockMekonnen/go-clean-starter/internal/user/app/events"
)


// MakeSendOTPEventConsumer subscribes to the OTP.SendOTPEvent and handles it
func SendOTPEventListener(subscriber events.Subscriber) func() error {
	return events.MakeEventConsumer(
		subscriber,
		events.EventAddress{
			Topic:     "OTP",
			EventType: "SendOTPEvent",
		},
		func() func(events.EventInterface) error {
			return func(evt events.EventInterface) error {
				// Cast the generic interface to your specific type
				event, ok := evt.(events.Event[app.SendOTPEventPayload])
				if !ok {
					return fmt.Errorf("invalid event type: expected SendOTPEventPayload")
				}

				payload := event.Payload

				// Simulate processing (e.g., send email or SMS)
				fmt.Printf("📩 Sending OTP to %s for %s\n", payload.Email, payload.OTPFor)

				// TODO: use your real email/SMS service here via deps

				return nil
			}
		},
		nil, // or pass &events.SubscriberOptions{Single: true} if needed
	)
}
