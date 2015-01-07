package facade_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOsagent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Osagent Suite")
}
