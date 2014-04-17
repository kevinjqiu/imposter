package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
)

type Error struct {
	Message string `json:"message"`
}

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

type PresetList []Preset

type Presets map[string]PresetList

func (this Presets) Add(preset *Preset) {
	rule := rule(preset.Matcher.Method, preset.Matcher.Endpoint)
	presets_for_rule, ok := this[rule]
	if !ok {
		presets_for_rule = PresetList(make([]Preset, 0, 5))
	}
	presets_for_rule = append(presets_for_rule, *preset)
	this[rule] = presets_for_rule
}

var presets = Presets(make(map[string]PresetList))

func rule(method string, endpoint string) string {
	return fmt.Sprintf("%s %s", method, endpoint)
}

func GetPreset(enc encoder.Encoder, w http.ResponseWriter) (int, []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return http.StatusOK, encoder.Must(enc.Encode(&presets))
}

func CreatePreset(
	params martini.Params, w http.ResponseWriter,
	r *http.Request, enc encoder.Encoder) (int, []byte) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest,
			encoder.Must(enc.Encode(&Error{err.Error()}))
	}
	preset := &Preset{}
	err = json.Unmarshal(bytes, &preset)
	if err != nil {
		return http.StatusBadRequest,
			encoder.Must(enc.Encode(&Error{err.Error()}))
	}
	presets.Add(preset)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return http.StatusCreated, encoder.Must(enc.Encode(preset))
}

func PresetRouter(r martini.Router) {
	r.Get("/", GetPreset)
	r.Post("/", CreatePreset)
}

func GetMock(req *http.Request, writer http.ResponseWriter, params martini.Params) (int, string) {
	method := req.Method
	endpoint := "/" + params["_1"]
	presets, ok := presets[rule(method, endpoint)]
	if !ok {
		return http.StatusNotFound, ""
	}

	for _, preset := range presets {
		if preset.Matcher.Match(req) {
			for key, value := range preset.Response.Headers {
				writer.Header().Set(key, value)
			}
			return preset.Response.StatusCode, preset.Response.Body
		}
	}
	return 418, "Not a teapot :("
}

func MockRouter(r martini.Router) {
	r.Get("/**", GetMock)
	r.Post("/**", GetMock)
}

func main() {
	m := martini.Classic()

	m.Use(func(c martini.Context) {
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
	})

	m.Group("/p", PresetRouter)
	m.Group("/m", MockRouter)

	m.Run()
}
