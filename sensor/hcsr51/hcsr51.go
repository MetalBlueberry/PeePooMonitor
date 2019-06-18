package hcsr51

//Automatically generate mocks
//go:generate mockgen -destination=../mocks/mock_$GOFILE -package=mocks  --source=$GOFILE

import (
	"fmt"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type HCSR51 struct {
	PinNumber     int
	NotifyTimeout time.Duration
	commander     Commander
}

func NewHCSR51(PinNumber int) *HCSR51 {
	return &HCSR51{
		PinNumber:     PinNumber,
		NotifyTimeout: time.Second * 1,
		commander:     Command{},
	}
}
func NewHCSR51Timeout(PinNumber int, NotifyTimeout time.Duration) *HCSR51 {
	return &HCSR51{
		PinNumber:     PinNumber,
		NotifyTimeout: NotifyTimeout,
		commander:     Command{},
	}
}

func (w *HCSR51) SetCommander(commander Commander) {
	w.commander = commander
}

type Commander interface {
	Command(command string, args ...string) ([]byte, error)
}
type Command struct{}

func (c Command) Command(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}

func (w *HCSR51) DetectMotion() <-chan int {
	log.Info("Using wiringpi library")
	notify := make(chan int)

	go func(w *HCSR51, notify chan<- int) {
		for {
			status, err := w.WatchInputChanges()
			if err != nil {
				log.WithError(err).Fatal("Unable to interact with GPIO")
			}
			select {
			case notify <- status:
			case <-time.After(w.NotifyTimeout):
				log.Error("There is no listener for motion detection notifications")
				close(notify)
				return
			}
		}
	}(w, notify)

	return notify
}

func (w *HCSR51) WatchInputChanges() (int, error) {
	defer func() {
		if r := recover(); r!=nil{
		log.WithField("recover", r).Error("Panic recoverin Detect Motion, Should only happen during test conditions")
		debug.PrintStack()
		}
	}()

	log.Debug("Waiting for edge")
	command := "gpio"
	args := fmt.Sprintf("-g wfi %d both", w.PinNumber)
	_, err := w.commander.Command(command, strings.Split(args, " ")...)
	if err != nil {
		return -1, err
	}
	log.Debug("detected edge")

	return w.Status()
}

func (w *HCSR51) Status() (int, error) {
	command := "gpio"
	args := fmt.Sprintf("-g read %d", w.PinNumber)
	out, err := w.commander.Command(command, strings.Split(args, " ")...)
	if err != nil {
		return -1, err
	}
	log.WithField("state", out).Debug("State readed")
	value, err := strconv.Atoi(strings.Trim(string(out), " \n"))
	if err != nil {
		return -1, err
	}
	return value, nil
}