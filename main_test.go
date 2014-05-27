package imposter_test

import (
	"fmt"
	"regexp"
	. "github.com/kevinjqiu/imposter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func BeSomeUUID() OmegaMatcher {
	return &uuidMatcher{}
}

type uuidMatcher struct {
}

func (matcher *uuidMatcher) Match(actual interface{}) (success bool, err error) {
	value, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("expects an uuid as string")
	}
	re := regexp.MustCompile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")
	return re.MatchString(value), nil
}

func (matcher *uuidMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nTo Be a UUID string", actual)
}

func (matcher *uuidMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nNot To Be a UUID string", actual)
}

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

				Expect(rule.Id).Should(BeSomeUUID())
			})
		})
	})
})
