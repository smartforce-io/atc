package apiserver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type maskResponseWriter struct {
	header     http.Header
	status     string
	statusCode int
}

func (maskResponseWriter *maskResponseWriter) Header() http.Header {
	return maskResponseWriter.header
}

func (maskResponseWriter *maskResponseWriter) Write(bytes []byte) (int, error) {
	maskResponseWriter.status = string(bytes)
	return 200, nil
}

func (maskResponseWriter *maskResponseWriter) WriteHeader(statusCode int) {
	maskResponseWriter.statusCode = statusCode
}

func TestWebhook(t *testing.T) {
	act := &ActApiServer{}

	var tests = []struct {
		event              string
		strForBody         string
		expectedStatus     string
		expectedStatusCode int
	}{
		{"delete", "", "success", 200},
		{"create", "", "success", 200},
		{"marketplace_purchase", "marketplace_body", "", 200},
		{"push", `{"installation": {"id": 8}}`, "", 200},
		{"def", "", "This webhook is undefined yet.", 404},
		{"", "", "This webhook is undefined yet.", 404},
	}

	for _, test := range tests {
		var jsonStr = []byte(test.strForBody)
		req := &http.Request{
			Body:   ioutil.NopCloser(bytes.NewBufferString(string(jsonStr))),
			Header: make(http.Header),
		}
		resp := &maskResponseWriter{}
		req.Header.Set("X-GitHub-Event", test.event)
		act.webhook(resp, req)
		if resp.status != test.expectedStatus {
			t.Errorf("wrong webhook! status error event: %q; expected:%q, got:%q", test.event, test.expectedStatus, resp.status)
		}
		if resp.statusCode != test.expectedStatusCode {
			t.Errorf("wrong webhook! statusCode error event: %q; expected:%d, got:%d", test.event, test.expectedStatusCode, resp.statusCode)
		}

	}

}
