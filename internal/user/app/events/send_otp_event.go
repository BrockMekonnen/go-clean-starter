package events

import (
	"github.com/BrockMekonnen/go-clean-starter/core/lib/events"
	"github.com/google/uuid"
)

type OTPPurpose string

const (
	Verification OTPPurpose = "Verification"
	TwoFA        OTPPurpose = "2FA"
	Forget       OTPPurpose = "Forget"
)

type SendOTPEventPayload struct {
	Email  string     `json:"email"`
	OTPFor OTPPurpose `json:"otpFor"`
}

func NewSendOTPEvent(email string, otpFor OTPPurpose) events.Event[SendOTPEventPayload] {
	return events.Event[SendOTPEventPayload]{
		EventAddress: events.EventAddress{
			Topic:     "OTP",
			EventType: "SendOTPEvent",
		},
		EventID: uuid.New().String(),
		Payload: SendOTPEventPayload{
			Email:  email,
			OTPFor: otpFor,
		},
	}
}
