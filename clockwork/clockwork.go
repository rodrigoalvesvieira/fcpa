package clockwork

import (
	ts "fcpa/timeseries"
	"time"
)

// Start starts a cycle that runs a function
func Start(fn func(ts.Result, ts.Result), usdtiers ts.Result, trytiers ts.Result) {
	cycletime := time.Second * 60 / 10
	ticker := time.NewTicker(cycletime)

	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(time.Second * 60 * 60 * 24 * 14)
		done <- true
	}()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fn(usdtiers, trytiers)
		}
	}
}
