package notification

// Notifier is an simple abstraction of different notifiers(sms, email, slack etc..)
type Notifier interface {
	// Notify push the content to the `receipient`.
	// receipient can phone number, email based type of notifier.
	Notify(content []byte, recipient string) error
}
