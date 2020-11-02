package signal

import (
	"context"
	"reflect"
	"time"
)

// After returns a read only channel that closes after the given duration
func After(duration time.Duration) <-chan struct{} {
	ret := make(chan struct{})
	go func() {
		<-time.After(duration)
		close(ret)
	}()
	return ret
}

// Combine return a read only channel that closes after receiving anything from one of
// the given channels
func Combine(signals ...<-chan struct{}) <-chan struct{} {
	ret := make(chan struct{})
	if len(signals) == 0 {
		close(ret)
		return ret
	}
	cases := make([]reflect.SelectCase, len(signals))
	for i, ch := range signals {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	go func() {
		_, _, ok := reflect.Select(cases)
		if !ok {
			close(ret)
		} else {
			ret <- struct{}{}
		}
	}()
	return ret
}

// Context returns a context that is cancelled upon the given signal
func Context(signal <-chan struct{}) context.Context {
	ret, cancel := context.WithCancel(context.Background())
	go func() {
		<-signal
		cancel()
	}()
	return ret
}
