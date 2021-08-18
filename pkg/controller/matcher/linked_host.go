package matcher

type linkedHost struct {
	pre  *linkedHost
	next *linkedHost
	host *Host
}

func newLinkedHost(host *Host) *linkedHost {
	return &linkedHost{
		host: host,
	}
}

func (l *linkedHost) append(host *Host) {
	if l.isLast() {
		node := newLinkedHost(host)
		node.pre = l
		l.next = node
	} else {
		l.next.append(host)
	}
}

func (l *linkedHost) remove() int {
	if l.isRoot() {
		return -1
	}

	if l.isLast() {
		l.pre.next = nil
		return 0
	}

	l.pre.next = l.next
	return 0
}

func (l *linkedHost) isRoot() bool {
	return l.pre == nil
}

func (l *linkedHost) isLast() bool {
	return l.next == nil
}

func (l *linkedHost) find(host Host) *linkedHost {
	if l.host.Equals(host) {
		return l
	}

	if l.next != nil {
		return l.next.find(host)
	}

	return nil
}
