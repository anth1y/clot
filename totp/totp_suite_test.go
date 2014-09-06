package totp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTotp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Totp Suite")
}
