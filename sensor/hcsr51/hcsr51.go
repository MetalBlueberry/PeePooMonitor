package hcsr51

//Automatically generate mocks
//go:generate mockgen -destination=../mocks/mock_$GOFILE -package=mocks  --source=$GOFILE

import (
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"

	// layout https://webofthings.org/wp-content/uploads/2016/10/pi-gpio.png
	"periph.io/x/periph/host/rpi"

	log "github.com/sirupsen/logrus"
)

type HCSR51 struct {
	Pin           gpio.PinIO
	NotifyTimeout time.Duration
}

func (h *HCSR51) InitController() {
	state, err := host.Init()
	if err != nil {
		log.WithError(err).Panic("When initializing controller")
	}

	log.Debug(state)

	if !rpi.Present() {
		log.Println("Not runing in Raspberry PI, Attach virtual pin")
		vpin := &VirtualPin{
			EdgeDuration: time.Second * 1,
			EdgePeriod:   time.Second * 2,
		}
		go vpin.Simulate(time.After(time.Second * 30))
		h.Pin = vpin
	}

	h.Pin.In(gpio.Float, gpio.BothEdges)
}

func NewHCSR51(pin gpio.PinIO) *HCSR51 {
	return NewHCSR51Timeout(pin, time.Second*1)
}

func NewHCSR51Timeout(pin gpio.PinIO, NotifyTimeout time.Duration) *HCSR51 {
	sensor := &HCSR51{
		Pin:           pin,
		NotifyTimeout: NotifyTimeout,
	}
	sensor.InitController()
	return sensor
}

func (w *HCSR51) DetectMotion(done <-chan struct{}) <-chan gpio.Level {
	log.Info("Using periph library")
	notify := make(chan gpio.Level)

	go func(w *HCSR51, notify chan<- gpio.Level, done <-chan struct{}) {
		edgeDetected := make(chan struct{})
		go w.WaitForEdgeChannel(edgeDetected)
		for {
			log.Debug("Waiting for edge")
			select {
			case <-edgeDetected:
				log.Debug("detected edge")
				status := w.Pin.Read()
				log.WithField("state", status).Debug("State readed")

				select {
				case notify <- status:
				case <-time.After(w.NotifyTimeout):
					log.Error("There is no listener for motion detection notifications")
					close(notify)
					return
				}
			case <-done:
				log.Info("Stop DetectMotion as done is requested")
				close(notify)
				return
			}

		}
	}(w, notify, done)

	return notify
}

func (h *HCSR51) WaitForEdgeChannel(edgeDetected chan<- struct{}) {
	for {
		edge := h.Pin.WaitForEdge(time.Second * 10)
		if edge {
			edgeDetected <- struct{}{}
		}
	}
}

func (h *HCSR51) Status() gpio.Level {
	return h.Pin.Read()
}
