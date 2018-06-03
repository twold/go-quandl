package endpoints_test

import (
	. "github.com/twold/go-quandl/endpoints"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("New", func() {
	Context("When I input the name of an invalid service", func() {
		It("returns the endpoint data", func() {
			_, err := New("dataset", "WIKI", "FB", "data", "json")
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("When I input the name of a valid service", func() {
		It("returns the endpoint data", func() {
			actual, err := New("datasets", "WIKI", "FB", "data", "json")
			Expect(err).Should(BeNil())
			Expect(actual.URL).Should(Equal("https://www.quandl.com/api/v3/datasets/WIKI/FB/data.json"))

		})
	})

	Context("When I input the name of a valid service", func() {
		It("returns the endpoint data", func() {
			actual, err := New("datasets", "WIKI", "GE", "data", "json")
			Expect(err).Should(BeNil())
			Expect(actual.URL).Should(Equal("https://www.quandl.com/api/v3/datasets/WIKI/GE/data.json"))

		})
	})

	Context("When I input the name of a valid service", func() {
		It("returns the endpoint data", func() {
			actual, err := New("datasets", "EOD", "FB", "metadata", "xml")
			Expect(err).Should(BeNil())
			Expect(actual.URL).Should(Equal("https://www.quandl.com/api/v3/datasets/EOD/FB/metadata.xml"))

		})
	})
})
