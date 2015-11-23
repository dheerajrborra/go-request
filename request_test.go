package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/blendlabs/go-assert"
)

type testObject struct {
	Id           int       `json:"id" xml:"id"`
	Name         string    `json:"name" xml:"name"`
	TimestampUtc time.Time `json:"timestamp_utc" xml:"timestamp_utc"`
	Value        float64   `json:"value" xml:"value"`
}

func newTestObject() testObject {
	to := testObject{}
	to.Id = rand.Int()
	to.Name = fmt.Sprintf("Test Object %d", to.Id)
	to.TimestampUtc = time.Now().UTC()
	to.Value = rand.Float64()
	return to
}

func okMeta() *HttpResponseMeta {
	return &HttpResponseMeta{StatusCode: http.StatusOK}
}

func errorMeta() *HttpResponseMeta {
	return &HttpResponseMeta{StatusCode: http.StatusInternalServerError}
}

func notFoundMeta() *HttpResponseMeta {
	return &HttpResponseMeta{StatusCode: http.StatusNotFound}
}

func writeJson(w http.ResponseWriter, meta *HttpResponseMeta, response interface{}) error {
	bytes, err := json.Marshal(response)
	if err == nil {
		if !isEmpty(meta.ContentType) {
			w.Header().Set("Content-Type", meta.ContentType)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		for key, value := range meta.Headers {
			w.Header().Set(key, strings.Join(value, ";"))
		}

		w.WriteHeader(meta.StatusCode)
		count, write_error := w.Write(bytes)
		if count == 0 {
			return errors.New("WriteJson : Didnt write any bytes.")
		}
		if write_error != nil {
			return write_error
		}
	} else {
		return err
	}
	return nil
}

func mockEchoEndpoint(meta *HttpResponseMeta) *httptest.Server {
	return getMockServer(func(w http.ResponseWriter, r *http.Request) {
		if !isEmpty(meta.ContentType) {
			w.Header().Set("Content-Type", meta.ContentType)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		for key, value := range meta.Headers {
			w.Header().Set(key, strings.Join(value, ";"))
		}

		defer r.Body.Close()
		bytes, _ := ioutil.ReadAll(r.Body)
		w.Write(bytes)
	})
}

func mockEndpoint(meta *HttpResponseMeta, returnWithObject interface{}) *httptest.Server {
	return getMockServer(func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, meta, returnWithObject)
	})
}

func getMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestCreateHttpRequestWithUrl(t *testing.T) {
	assert := assert.New(t)
	sr := NewRequest().
		WithUrl("http://localhost:5001/api/v1/path/2?env=dev&foo=bar")

	assert.Equal("http", sr.Scheme)
	assert.Equal("localhost:5001", sr.Host)
	assert.Equal("GET", sr.Verb)
	assert.Equal("/api/v1/path/2", sr.Path)
	assert.Equal([]string{"dev"}, sr.QueryString["env"])
	assert.Equal([]string{"bar"}, sr.QueryString["foo"])
	assert.Equal(2, len(sr.QueryString))
}

func TestHttpGet(t *testing.T) {
	assert := assert.New(t)
	returned_object := newTestObject()
	ts := mockEndpoint(okMeta(), returned_object)
	test_object := testObject{}
	meta, err := NewRequest().AsGet().WithUrl(ts.URL).FetchJsonToObjectWithMeta(&test_object)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returned_object, test_object)
}

func TestHttpPostWithJsonBody(t *testing.T) {
	assert := assert.New(t)

	returned_object := newTestObject()
	ts := mockEchoEndpoint(okMeta())

	test_object := testObject{}
	meta, err := NewRequest().AsPost().WithUrl(ts.URL).WithJsonBody(&returned_object).FetchJsonToObjectWithMeta(&test_object)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returned_object, test_object)
}

func TestHttpPostWithXmlBody(t *testing.T) {
	assert := assert.New(t)

	returned_object := newTestObject()
	ts := mockEchoEndpoint(okMeta())

	test_object := testObject{}
	meta, err := NewRequest().AsPost().WithUrl(ts.URL).WithXmlBody(&returned_object).FetchXmlToObjectWithMeta(&test_object)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(returned_object, test_object)
}
