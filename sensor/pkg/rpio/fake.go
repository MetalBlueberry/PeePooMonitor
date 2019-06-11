// +build !arm

package rpio

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (r HCSR51) DetectMotion() <-chan bool {
	log.Info("Running in simulation mode")
	notify := make(chan bool)

	go func(notify chan<- bool) {
		var status = false
		for {
			notify <- status
			status = !status
			time.Sleep(time.Second * 5)
		}
	}(notify)

	return notify

}
