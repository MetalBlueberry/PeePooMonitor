package hcsr51_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHcsr51(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hcsr51 Suite")
}
