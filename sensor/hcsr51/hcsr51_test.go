package hcsr51_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"
	. "github.com/metalblueberry/PeePooMonitor/sensor/hcsr51"
	"github.com/metalblueberry/PeePooMonitor/sensor/mocks"
)

var _ = Describe("Given the sensor HCSR51", func() {
	var (
		mockCtrl *gomock.Controller //gomock struct
		// generated using mockgen command
		mockCommander *mocks.MockCommander
		sensor        *HCSR51
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCommander = mocks.NewMockCommander(mockCtrl)
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("for the function WatchInputChanges", func() {
		BeforeEach(func() {
			sensor = NewHCSR51Timeout(17, time.Millisecond*10)
			sensor.SetCommander(mockCommander)
		})
		Describe("When the status changes to Enabled", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return([]byte("\n"), nil).Times(1),
					mockCommander.EXPECT().Command("gpio", []string{"-g", "read", "17"}).Return([]byte("1\n"), nil).Times(1),
				)
			})
			It("should return the device status True", func() {
				got, err := sensor.WatchInputChanges()
				Expect(err).To(BeNil())
				Expect(got).To(BeTrue())
			})
		})
		Describe("When the status changes to Disabled", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return([]byte("\n"), nil).Times(1),
					mockCommander.EXPECT().Command("gpio", []string{"-g", "read", "17"}).Return([]byte("0\n"), nil).Times(1),
				)
			})
			It("should return the device status False", func() {
				got, err := sensor.WatchInputChanges()
				Expect(err).To(BeNil())
				Expect(got).To(BeFalse())
			})
		})

		Describe("When the command execution fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return(nil, errors.New("execution error")).Times(1),
				)
			})
			It("Should return an error", func() {
				_, err := sensor.WatchInputChanges()
				Expect(err).ToNot(BeNil())
			})
		})
		Describe("When the command output is not as expected", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return([]byte("\n"), nil).Times(1),
					mockCommander.EXPECT().Command("gpio", []string{"-g", "read", "17"}).Return([]byte("Wrong output\n"), nil).Times(1),
				)
			})
			It("Should return and error", func() {
				_, err := sensor.WatchInputChanges()
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("For the function DetectMotion", func() {
		BeforeEach(func() {
			sensor = NewHCSR51Timeout(17, time.Millisecond*10)
			sensor.SetCommander(mockCommander)
		})
		Describe("When the status changes to Enabled", func() {
			BeforeEach(func() {
				mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return([]byte("\n"), nil).AnyTimes()
				mockCommander.EXPECT().Command("gpio", []string{"-g", "read", "17"}).Return([]byte("1\n"), nil).AnyTimes()
			})
			It("Shoud be notified", func() {
				notifications := sensor.DetectMotion()
				event := <-notifications
				Expect(event).To(BeTrue())
			})
		})
		Describe("When the status changes to Disabled", func() {
			BeforeEach(func() {
				mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return([]byte("\n"), nil).AnyTimes()
				mockCommander.EXPECT().Command("gpio", []string{"-g", "read", "17"}).Return([]byte("0\n"), nil).AnyTimes()
			})
			It("Shoud be notified", func() {
				notifications := sensor.DetectMotion()
				event := <-notifications
				Expect(event).To(BeFalse())
			})
		})
		Describe("When the there is no reader in notifications channel", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockCommander.EXPECT().Command("gpio", []string{"-g", "wfi", "17", "both"}).Return([]byte("\n"), nil).Times(1),
					mockCommander.EXPECT().Command("gpio", []string{"-g", "read", "17"}).Return([]byte("0\n"), nil).Times(1),
				)
			})
			It("Should be closed", func() {
				notifications := sensor.DetectMotion()
				Expect(notifications).ToNot(BeClosed())
				time.Sleep(sensor.NotifyTimeout + time.Millisecond + 10)
				Expect(notifications).To(BeClosed())
			})
		})
	})
})
