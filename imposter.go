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

type Preset struct {
	Method     string                 `json:"method"`
	Endpoint   string                 `json:"endpoint"`
	Body       string                 `json:"body"`
	Header     map[string]interface{} `json:"header"`
	StatusCode int                    `json:"status_code"`
}

var presets = make(map[string]Preset)

func spec(method string, endpoint string) string {
	return fmt.Sprintf("%s %s", method, endpoint)
}

func GetPreset(enc encoder.Encoder) (int, []byte) {
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
	presets[spec(preset.Method, preset.Endpoint)] = *preset
	return http.StatusCreated, encoder.Must(enc.Encode(preset))
}

func PresetRouter(r martini.Router) {
	r.Get("/", GetPreset)
	r.Post("/", CreatePreset)
}

func GetMock(req *http.Request, params martini.Params) (int, string) {
	method := req.Method
	endpoint := "/" + params["_1"]
	preset, ok := presets[spec(method, endpoint)]
	if !ok {
		return http.StatusNotFound, "{}"
	}
	return preset.StatusCode, preset.Body
}

func MockRouter(r martini.Router) {
	r.Get("/**", GetMock)
	r.Post("/**", GetMock)
}

func main() {
	m := martini.Classic()

	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	m.Group("/p", PresetRouter)
	m.Group("/m", MockRouter)

	m.Run()
}
