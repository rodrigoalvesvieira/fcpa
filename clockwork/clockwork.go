package clockwork

import (
	"time"
	ts "operatordesk/timeseries"
)


func Start(fn func(ts.Result, ts.Result), usdtiers ts.Result, trytiers ts.Result) {
	cycle_time := time.Second * 60 / 12 //* 2
	ticker := time.NewTicker(cycle_time)

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
