package totp_test

import (
	"time"
	. "totp/totp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Totp", func() {
	token := Token{
		Id:       1,
		Label:    "username@service",
		Secret:   "ptcuyfm2sjjqgh5c",
		Digest:   "sha1",
		Digits:   6,
		TimeStep: 30,
	}

	ts, err := time.Parse("January 2 15:04:05 UTC 2006", "September 6 17:16:28 UTC 2014")
	if err != nil {
		panic(err)
	}

	Describe("Now()", func() {
		It("should return the OTP for the current time", func() {
			Expect(len(token.Now())).To(Equal(6))
		})
	})

	Describe("At(time)", func() {
		It("should return the OTP for a given time", func() {
			Expect(token.At(ts)).To(Equal("598299"))
		})
	})

	Describe("DecodeSecret()", func() {
		It("Decodes the Base32 secret", func() {
			key := token.Key()

			Expect(key).To(Equal([]byte{
				124,
				197,
				76,
				21,
				154,
				146,
				83,
				3,
				31,
				162,
			}))
		})
	})
})
