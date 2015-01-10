package jobstore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJobstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jobstore Suite")
}
