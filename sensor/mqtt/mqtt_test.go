package mqtt_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	//. "github.com/onsi/gomega"

	"github.com/metalblueberry/PeePooMonitor/sensor/mocks"
	. "github.com/metalblueberry/PeePooMonitor/sensor/mqtt"
)

var _ = Describe("Mqtt", func() {
	var (
		mockCtrl *gomock.Controller //gomock struct
		// generated using mockgen command
		mockMqttPublisher *mocks.MockmqttPublisher
		// sensor        *HCSR51
		token  *mocks.MockToken
		client *MqttClient
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockMqttPublisher = mocks.NewMockmqttPublisher(mockCtrl)

		token = mocks.NewMockToken(mockCtrl)
		opts := &MqttClientOptions{
			Server:      "server",
			Qos:         1,
			Clientid:    "testClient",
			Username:    "Username",
			Password:    "Password",
			SendTimeout: 1,
			OnConnect: func(client Publisher) {
			},
		}
		client = &MqttClient{
			Options: opts,
			Client:  mockMqttPublisher,
		}
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})
	Describe("When a Power Status update is published", func() {
		Describe("Is false", func() {
			It("Should forwared PowerOff", func() {
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				mockMqttPublisher.EXPECT().Publish(
					gomock.Eq(string(PowerStatusAddress)),
					gomock.Any(),
					true,
					"PowerOff",
				).Return(token).Times(1)

				client.PublishPowerStatus(false)
			})
		})
		Describe("Is true", func() {
			It("Should forwared PowerOn", func() {
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				mockMqttPublisher.EXPECT().Publish(
					gomock.Eq(string(PowerStatusAddress)),
					gomock.Any(),
					true,
					"PowerOn",
				).Return(token).Times(1)

				client.PublishPowerStatus(true)
			})
		})
	})
	Describe("When a Sensor Status update is published", func() {
		It("Should forwared the same message", func() {
			status := "TestStatus"
			token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
			mockMqttPublisher.EXPECT().Publish(
				gomock.Eq(string(SensorStatusAddress)),
				gomock.Any(),
				true,
				status,
			).Return(token).Times(1)

			client.PublishSensorStatus(status)
		})
	})
})
