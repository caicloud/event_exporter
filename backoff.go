package main

import (
	"sync"
	"time"

	"k8s.io/utils/integer"
)

type backoffEntry struct {
	backoff    time.Duration
	lastUpdate time.Time
}

// Backoff stores entries hoding backoff infomation
type Backoff struct {
	sync.Mutex
	baseDuration time.Duration
	maxDuration  time.Duration
	perItemEntry map[string]*backoffEntry
}

// NewBackoff returns a brand new backoff storage
func NewBackoff(init, max time.Duration) *Backoff {
	backoff := &Backoff{
		baseDuration: init,
		maxDuration:  max,
		perItemEntry: map[string]*backoffEntry{},
	}
	return backoff
}

// Next moves backoff to next mark, capping at maxDuration
func (p *Backoff) Next(id string, count int, eventTime time.Time) {
	p.Lock()
	defer p.Unlock()
	entry, ok := p.perItemEntry[id]
	if !ok || hasExpired(time.Now(), entry.lastUpdate, p.maxDuration) {
		p.perItemEntry[id] = &backoffEntry{
			backoff:    p.baseDuration,
			lastUpdate: eventTime,
		}
	} else {
		delay := entry.backoff * 2
		entry.lastUpdate = eventTime
		entry.backoff = time.Duration(integer.Int64Min(int64(delay), int64(p.maxDuration)))
	}
}

// Reset forces clearing of all backoff data for a given key.
func (p *Backoff) Reset(id string) {
	p.Lock()
	defer p.Unlock()
	delete(p.perItemEntry, id)
}

// IsInBackOffSinceUpdate returns True if time since lastupdate is less than the current backoff window.
func (p *Backoff) IsInBackOffSinceUpdate(id string, eventTime time.Time) bool {
	p.Lock()
	defer p.Unlock()
	entry, ok := p.perItemEntry[id]
	if !ok {
		return false
	}
	if hasExpired(eventTime, entry.lastUpdate, p.maxDuration) {
		return false
	}
	return eventTime.Sub(entry.lastUpdate) < entry.backoff
}

// After 2*maxDuration we restart the backoff factor to the beginning
func hasExpired(eventTime time.Time, lastUpdate time.Time, maxDuration time.Duration) bool {
	return eventTime.Sub(lastUpdate) > maxDuration*2 // consider stable if it's ok for twice the maxDuration
}

// GC records that have aged past maxDuration. Backoff users are expected
// to invoke this periodically.
func (p *Backoff) GC() {
	p.Lock()
	defer p.Unlock()
	now := time.Now()
	for id, entry := range p.perItemEntry {
		if now.Sub(entry.lastUpdate) > p.maxDuration*2 {
			// GC when entry has not been updated for 2*maxDuration
			delete(p.perItemEntry, id)
		}
	}
}

// Get return the backoff entry content for specific id
func (p *Backoff) Get(id string) (time.Duration, time.Time) {
	return p.perItemEntry[id].backoff, p.perItemEntry[id].lastUpdate
}

// AllKeysStateSinceUpdate returns all keys with the state of whether in backoff
func (p *Backoff) AllKeysStateSinceUpdate(eventTime time.Time) map[string]bool {
	p.Lock()
	defer p.Unlock()
	state := map[string]bool{}
	for id, entry := range p.perItemEntry {
		if hasExpired(eventTime, entry.lastUpdate, p.maxDuration) {
			continue
		} else if eventTime.Sub(entry.lastUpdate) < entry.backoff {
			state[id] = true
		} else if eventTime.Sub(entry.lastUpdate) >= entry.backoff {
			state[id] = false
		}
	}
	return state
}
