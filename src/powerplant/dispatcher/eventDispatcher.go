package dispatcher

type EventDispatcher struct {
	listeners map[string][]func(interface{})
}

type IDispatcher interface {
	AddListener(name string, callback func(interface{}))
}

func NewEventDispatcher() *EventDispatcher {
	dispatcher := EventDispatcher{
		listeners: make(map[string][]func(interface{})),
	}
	return &dispatcher
}

func (ed *EventDispatcher) AddListener(name string, callback func(interface{})) {
	ed.listeners[name] = append(ed.listeners[name], callback)
}

func (ed *EventDispatcher) PublishEvent(name string, data interface{}) {
	if ed.listeners[name] != nil {
		for _, callback := range ed.listeners[name] {
			callback(data)
		}
	}
}
