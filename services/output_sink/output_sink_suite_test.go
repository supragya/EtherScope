package outputsink_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOutputSink(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OutputSink Suite")
}
