package notifications

// MultiNotifier fans an event out to every wrapped notifier in order.
//
// MultiNotifier is safe for concurrent use as long as every wrapped notifier
// is. Each delegate is responsible for its own goroutine semantics; the
// MultiNotifier itself does no goroutine work — it just iterates and calls
// NotifyAsync on each.
type MultiNotifier []Notifier

// NotifyAsync forwards the event to every underlying notifier. nil entries
// are skipped so callers can compose the slice conditionally without
// guarding each element.
func (m MultiNotifier) NotifyAsync(event *BookingEvent) {
	for _, n := range m {
		if n == nil {
			continue
		}
		n.NotifyAsync(event)
	}
}
