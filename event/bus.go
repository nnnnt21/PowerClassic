package event

type Event interface {
	IsCancelled() bool
	Cancel()
}

type BaseEvent struct {
	cancelled bool
}

func (e *BaseEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *BaseEvent) Cancel() {
	e.cancelled = true
}

type EventHandler func(Event)

type EventBus struct {
	handlers map[string][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (bus *EventBus) Register(eventName string, handler EventHandler) {
	bus.handlers[eventName] = append(bus.handlers[eventName], handler)
}

func (bus *EventBus) Fire(eventName string, event Event) {
	if handlers, ok := bus.handlers[eventName]; ok {
		for _, handler := range handlers {
			handler(event)
			if event.IsCancelled() {
				break
			}
		}
	}
}
