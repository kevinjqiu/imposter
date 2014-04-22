package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type PresetRequestMatcher struct {
	Method   string            `json:"method"`
	Endpoint string            `json:"endpoint"`
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
}

func (this *PresetRequestMatcher) Match(req *http.Request) bool {
	fmt.Printf("%s", req.Header)
	// match header
	for key, value := range this.Headers {
		headerValues, ok := req.Header[key]
		if !ok {
			return false
		}
		for _, requestValue := range headerValues {
			if requestValue != value {
				return false
			}
		}
	}
	// match body
	bytes, _ := ioutil.ReadAll(req.Body)
	if string(bytes) != this.Body {
		return false
	}
	return true
}

type PresetResponse struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"status_code"`
}

type Preset struct {
	Matcher  PresetRequestMatcher `json:"matcher"`
	Response PresetResponse       `json:"response"`
}
