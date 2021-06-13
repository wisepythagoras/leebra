package event_target

// Event describes the events that are dispatched to event listeners.
type Event struct {
	Type       string
	TimeStamp  int64
	Bubbles    interface{}
	Cancelable bool
	detail     interface{}
}
