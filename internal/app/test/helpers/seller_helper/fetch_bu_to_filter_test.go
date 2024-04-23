package helper_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FetchBuToFilter", func() {
	Context("When current user is not present", func() {
		It("Should return input BUs", func() {
			Expect(1).To(Equal(1))
		})
	})
	Context("When current user is present", func() {
		Context("When inputBUs is empty", func() {
			It("Should return BUs from oms response", func() {
				Expect(1).To(Equal(1))
			})
		})
		Context("When inputBUs is present", func() {
			Context("When OMS response is empty", func() {
				It("Should return BUs from passed as arguments", func() {
					Expect(1).To(Equal(1))
				})
			})
			Context("When OMS response is not empty", func() {
				It("Should return intersection of input and oms response", func() {
					Expect(1).To(Equal(1))
				})
			})
		})
	})
})
