package imposter

import "code.google.com/p/go-uuid/uuid"

type Preset map[Rule]Response

func (preset Preset) AddRule(rule *Rule, response Response) {
	rule.Id = uuid.New()
	preset[*rule] = response
}

type Rule struct {
	Id     string
	Path   string
	Method string
}

func (rule Rule) Match() bool {
	return false
}

type Response struct {
	Status int
	Header map[string]string
}
