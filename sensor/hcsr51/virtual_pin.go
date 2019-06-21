package hcsr51

import (
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
)

//Fake pin when runining in computer
type VirtualPin struct {
	EdgePeriod   time.Duration
	EdgeDuration time.Duration
	status       bool
	notify       chan struct{}
}

func (VirtualPin) String() string {
	panic("not implemented")
}

func (VirtualPin) Halt() error {
	panic("not implemented")
}

func (VirtualPin) Name() string {
	panic("not implemented")
}

func (VirtualPin) Number() int {
	panic("not implemented")
}

func (VirtualPin) Function() string {
	panic("not implemented")
}

func (v *VirtualPin) In(pull gpio.Pull, edge gpio.Edge) error {
	return nil
}

func (v *VirtualPin) Simulate(done <-chan time.Time) {
	for {
		time.Sleep(v.EdgePeriod)
		v.status = true
		if v.notify != nil {
			v.notify <- struct{}{}
		}

		time.Sleep(v.EdgeDuration)
		if v.notify != nil {
			v.notify <- struct{}{}
		}
		v.status = false

		select {
		case <-done:
			return
		default:
			continue
		}
	}
}

func (v *VirtualPin) Read() gpio.Level {
	return gpio.Level(v.status)
}

func (v *VirtualPin) WaitForEdge(timeout time.Duration) bool {
	v.notify = make(chan struct{})
	defer func() {
		close(v.notify)
		v.notify = nil
	}()
	select {
	case <-v.notify:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (VirtualPin) Pull() gpio.Pull {
	panic("not implemented")
}

func (VirtualPin) DefaultPull() gpio.Pull {
	panic("not implemented")
}

func (VirtualPin) Out(l gpio.Level) error {
	panic("not implemented")
}

func (VirtualPin) PWM(duty gpio.Duty, f physic.Frequency) error {
	panic("not implemented")
}
