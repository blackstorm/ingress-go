package controller

import "github.com/blackstorm/ingress-go/pkg/watcher"

type ChannelEvent struct {
	event  watcher.Event
	values []interface{}
}

func (ce ChannelEvent) GetFirst() interface{} {
	return ce.values[0]
}

func (ce ChannelEvent) GetLast() interface{} {
	return ce.values[len(ce.values)-1]
}

type WatcherEventChannel struct {
	ch chan ChannelEvent
}

func newWatcherEventChannel() WatcherEventChannel {
	return WatcherEventChannel{
		ch: make(chan ChannelEvent),
	}
}

func (ec WatcherEventChannel) send(channelEvent ChannelEvent) {
	ec.ch <- channelEvent
}

// TODO context cancel
func (ec WatcherEventChannel) listen(handler func(ChannelEvent)) {
	for {
		event := <-ec.ch
		handler(event)
	}
}
