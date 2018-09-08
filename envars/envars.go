package envars

import (
	"os"
)

// GetTwilioEnvVar returns the Twilio Api Key
func GetTwilioEnvVar() string {
	return os.Getenv("TWILIO_API_KEY")
}

// GetTwilioPhoneNumber returns the Twilio Phone Number that we use as *sender* for the SMS messages in this application.
func GetTwilioPhoneNumber() string {
	return os.Getenv("TWILIO_PHONE_NUMBER")
}

// GetKrakenAPIKey returns the Kraken API Key
func GetKrakenAPIKey() string {
	return os.Getenv("OPERATOR_DESK_KRAKEN_API_KEY")
}

// GetKrakenPrivateKey returns the Kraken Private Key
func GetKrakenPrivateKey() string {
	return os.Getenv("OPERATOR_DESK_KRAKEN_PRIVATE_KEY")
}
