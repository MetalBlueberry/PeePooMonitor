package main

import (
	"os"

	"github.com/metalblueberry/PeePooMonitor/pkg/rpio"
	"github.com/metalblueberry/PeePooMonitor/pkg/telegrambot"

	//"github.com/metalblueberry/PeePooMonitor/pkg/telegrambot"
	"time"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

var (
	pinNumber = 17
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	//log filename and line :D
	log.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

type motionSensor interface {
	DetectMotion() <-chan bool
}

type motionNotification interface {
	NotifyMotion() chan<- string
}

func main() {
	var err error
	TelegramToken, err := ioutil.ReadFile("/run/secrets/telegram_bot_token")
	if err != nil {
		log.WithError(err).Panic("cannot load telegram bot token from secrets")
	}

	sensor := rpio.HCSR51{
		PinNumber:         pinNumber,
		SamplingInSeconds: 1,
	}

	notifier := telegrambot.TelegramBot{
		Token: string(TelegramToken),
	}

	err = run(notifier ,sensor)
	if err != nil {
		log.WithError(err).Error("Program finished with error")
	}

}

func run(notifier motionNotification, sensor motionSensor) error {
	detectMotion := sensor.DetectMotion()
	notifyMotion := notifier.NotifyMotion()

	log.WithField("Pin Number", pinNumber).Info("Waiting for cats")

	presence := false
	for {
		select {
		case presence = <-detectMotion:
			if presence {
				log.Info("poop is comming!!")
				notifyMotion <- "Poop is comming"
			} else {
				log.Info("it's done")
				notifyMotion <- "It's done"
				
			}
		case <-time.After(60 * 60 * time.Second):
			log.WithField("presence", presence).Debug("last status")
		}
	}
}
