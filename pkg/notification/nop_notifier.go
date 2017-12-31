package notification

// nop is a no-op notifier. It notifies to nothing.
type nop struct {
}

// NewNop returns an instance of `nop` notifier.
func NewNop() Notifier {
	return &nop{}
}

// Notify implements the Notifier.
func (_ *nop) Notify(content []byte, to string) error {
	return nil
}
