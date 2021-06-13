package event_target

// EventTarget implements the
type EventTarget struct {
	ListenerMap map[string][]*CallbackFn
}

// EnsureInitialized makes sure that ListenerMap is not nil.
func (et *EventTarget) EnsureInitialized() {
	if et.ListenerMap == nil {
		et.ListenerMap = make(map[string][]*CallbackFn)
	}
}

// AddEventListener adds an event listener.
func (et *EventTarget) AddEventListener(t string, callback *CallbackFn) {
	et.EnsureInitialized()

	_, categoryExists := et.ListenerMap[t]

	if !categoryExists {
		et.ListenerMap[t] = []*CallbackFn{}
	}

	et.ListenerMap[t] = append(et.ListenerMap[t], callback)
}

// RemoveEventListener removes an event listener.
func (et *EventTarget) RemoveEventListener(t string, targetCallback *CallbackFn) {
	et.EnsureInitialized()

	callbacks, categoryExists := et.ListenerMap[t]

	for i, callback := range callbacks {
		if categoryExists && targetCallback == callback {
			// Remove the callback from the
			et.ListenerMap[t] = append(callbacks[:i], callbacks[i+1:]...)
			break
		}
	}
}

// DispatchEvent dispatches an event.
func (et *EventTarget) DispatchEvent(event *Event) {
	et.EnsureInitialized()

	if event == nil {
		return
	}

	callbacks, categoryExists := et.ListenerMap[event.Type]

	if categoryExists {
		for _, callback := range callbacks {
			(*callback)(event)
		}
	}
}
