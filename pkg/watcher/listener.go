package watcher

type EventListener interface {
	Update(event Event, updates ...interface{})
}
