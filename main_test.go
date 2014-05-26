package imposter_test

import (
	. "github.com/kevinjqiu/imposter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Imposter", func() {
	Describe("Presets", func() {
		Context("Adding a rule the preset", func() {
			It("should have a unique id field", func() {
				p := Preset{}
				rule := Rule{
					Path:   "/foo/bar",
					Method: "GET",
				}
				response := Response{}
				p.AddRule(&rule, response)

				Expect(rule.Id).To(Equal(""))
			})
		})
	})
})
