package masterclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMasterclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Masterclient Suite")
}
