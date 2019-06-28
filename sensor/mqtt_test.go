package main_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	. "github.com/metalblueberry/PeePooMonitor/sensor"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination=./mqtt_mock_test.go -package=main_test  github.com/eclipse/paho.mqtt.golang Client,Token

var _ = Describe("Mqtt", func() {
	var (
		mockCtrl *gomock.Controller //gomock struct
		// generated using mockgen command
		mockClient *MockClient
		// sensor        *HCSR51
		token *MockToken
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockClient = NewMockClient(mockCtrl)

		token = NewMockToken(mockCtrl)
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})
	Describe("When a Power Status update is published", func() {
		Describe("Is false", func() {
			var result error
			It("Should forwared PowerOff", func() {
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				token.EXPECT().Error().Return(nil).Times(1)
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressPowerStatus)),
					gomock.Any(),
					true,
					"PowerOff",
				).Return(token).Times(1)
				result = PublishPowerStatus(mockClient, false)
			})
			It("Should return nil error", func() {
				Expect(result).To(BeNil())
			})
		})
		Describe("Is true", func() {
			var result error
			It("Should forwared PowerOn", func() {
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				token.EXPECT().Error().Return(nil).Times(1)
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressPowerStatus)),
					gomock.Any(),
					true,
					"PowerOn",
				).Return(token).Times(1)

				result = PublishPowerStatus(mockClient, true)
			})
			It("Should return nil error", func() {
				Expect(result).To(BeNil())
			})
		})
		Describe("The message is not sent", func() {
			It("Should return an error if there is a send timeout", func() {
				status := "TestStatus"
				token.EXPECT().WaitTimeout(gomock.Any()).Return(false).Times(1)
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressSensorStatus)),
					gomock.Any(),
					false,
					status,
				).Return(token).Times(1)

				result := PublishPowerStatus(mockClient, true)
				Expect(result).ToNot(BeNil())
			})
			It("Should return an error if token contains an error", func() {
				status := "TestStatus"
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				token.EXPECT().Error().Return(errors.New("Test Error sending")).AnyTimes()
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressSensorStatus)),
					gomock.Any(),
					false,
					status,
				).Return(token).Times(1)

				result := PublishPowerStatus(mockClient, true)
				Expect(result).ToNot(BeNil())
			})
	})
	Describe("When a Sensor Status update is published", func() {
		Describe("The message is sent successfully", func() {
			var result error
			It("Should forwared the same message", func() {
				status := "TestStatus"
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				token.EXPECT().Error().Return(nil).Times(1)
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressSensorStatus)),
					gomock.Any(),
					false,
					status,
				).Return(token).Times(1)

				result = PublishSensorStatus(mockClient, status)
			})
			It("Should return nil error", func() {
				Expect(result).To(BeNil())
			})
		})
		Describe("The message is not sent", func() {
			It("Should return an error if there is a send timeout", func() {
				status := "TestStatus"
				token.EXPECT().WaitTimeout(gomock.Any()).Return(false).Times(1)
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressSensorStatus)),
					gomock.Any(),
					false,
					status,
				).Return(token).Times(1)

				result := PublishSensorStatus(mockClient, status)
				Expect(result).ToNot(BeNil())
			})
			It("Should return an error if token contains an error", func() {
				status := "TestStatus"
				token.EXPECT().WaitTimeout(gomock.Any()).Return(true).Times(1)
				token.EXPECT().Error().Return(errors.New("Test Error sending")).AnyTimes()
				mockClient.EXPECT().Publish(
					gomock.Eq(string(AddressSensorStatus)),
					gomock.Any(),
					false,
					status,
				).Return(token).Times(1)

				result := PublishSensorStatus(mockClient, status)
				Expect(result).ToNot(BeNil())
			})
		})
	})
})
