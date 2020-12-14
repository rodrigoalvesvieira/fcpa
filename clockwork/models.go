package clockwork

import (
	"reflect"
	"time"

	"fcpa/logger"
)

type Trigger chan struct{}

type Action interface {
	DoAction()
}

// Clock is a simple struct do an generic action recurrently
// It is also possible to use an external trigger to activate the action
// using channels
type Clock struct {
	action            Action
	waitTrigger       bool
	ticker            *time.Ticker
	trigger           Trigger
	quitCh            Trigger
	completionChannel Trigger
}

func (clock *Clock) Start() {
	go func() {
		if !clock.waitTrigger {
			clock.doAction()
		}
		for {
			select {
			case <-clock.ticker.C:
				logger.Info().Log("clock", reflect.TypeOf(clock.action).String())
				if !clock.waitTrigger {
					clock.doAction()
				}
			case <-clock.trigger:
				logger.Info().Log("trigger", reflect.TypeOf(clock.action).String())
				clock.waitTrigger = false
				clock.doAction()
			case <-clock.quitCh:
				clock.ticker.Stop()
				return
			}
		}
	}()
}

func (clock *Clock) doAction() {
	clock.action.DoAction()
	if clock.completionChannel != nil {
		// produce on channel if not full
		select {
		case clock.completionChannel <- struct{}{}:
		default:
			logger.Info().Log("msg", "action completed, but completion channel is full")
		}
	}
}

func (clock *Clock) Stop() {
	clock.quitCh <- struct{}{}
}

func NewClockwork(action Action, intervalSeconds time.Duration, trigger Trigger, waitTrigger bool, completionChannel Trigger) *Clock {
	return &Clock{
		action:            action,
		ticker:            time.NewTicker(intervalSeconds * time.Second),
		trigger:           trigger,
		quitCh:            make(chan struct{}),
		waitTrigger:       waitTrigger,
		completionChannel: completionChannel,
	}
}
