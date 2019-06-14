// +build arm

package rpio

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (w HCSR51) DetectMotion() <-chan bool {
	log.Info("Using wiringpi library")
	notify := make(chan bool)

	go func(pinNumber int, notify chan<- bool) {
		for {
			log.Debug("Waiting for edge")
			var command string
			var args string
			var err error
			var out []byte

			command = "gpio"
			args = fmt.Sprintf("-g wfi %d both", pinNumber)
			_, err = exec.Command(command, strings.Split(args, " ")...).Output()
			if err != nil {
				log.WithError(err).Fatal("Unable to interact with GPIO")
			}
			log.Debug("detected edge")

			command = "gpio"
			args = fmt.Sprintf("-g read %d", pinNumber)
			out, err = exec.Command(command, strings.Split(args, " ")...).Output()
			if err != nil {
				log.Panic(err)
			}
			log.WithField("state", out).Debug("State readed")
			notify <- "0" != string(out[0])
		}
	}(w.PinNumber, notify)

	return notify
}
