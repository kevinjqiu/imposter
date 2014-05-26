package imposter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestImposter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Imposter Suite")
}
