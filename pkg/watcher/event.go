package watcher

type Event string

const (
	Add    Event = "add"
	Update Event = "update"
	Delete Event = "delete"
)
