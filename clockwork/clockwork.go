package clockwork

import "time"

func Start(fn func()) {
	cycle_time := time.Second * 60  //* 2
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
			fn()
		}
	}
}
