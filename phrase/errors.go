package phrase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	*http.Response

	ErrorRaw        json.RawMessage `json:"error"`
	MessagesRaw     json.RawMessage `json:"messages"`
	Message         string          `json:"-"`
	ValidationError errorMap        `json:"-"`
}

func ResponseError(resp *http.Response) *ErrorResponse {
	e := &ErrorResponse{Response: resp}
	e.populate()
	return e
}

type errorMap map[string][]string

func (r *ErrorResponse) populate() {
	data, err := ioutil.ReadAll(r.Response.Body)
	if err == nil && data != nil {
		if contentType := r.Response.Header.Get("Content-Type"); strings.Index(contentType, "application/json") != -1 {
			r.populateFields(data)
		} else {
			r.Message = string(data)
		}
	}
}

func (r *ErrorResponse) populateFields(data []byte) {
	var err error
	if err = json.Unmarshal(data, r); err != nil {
		return
	}
	r.populateMessagesRaw()
	r.populateErrorRaw()
}

func (r *ErrorResponse) populateMessagesRaw() {
	var err error
	raw := r.MessagesRaw
	var s string
	if err = json.Unmarshal(raw, &s); err == nil {
		r.Message = s
		return
	}
	var m map[string][]string
	if err = json.Unmarshal(raw, &m); err == nil {
		r.ValidationError = m
		return
	}
	var a []string
	if err = json.Unmarshal(raw, &a); err == nil {
		r.ValidationError = map[string][]string{
			"error": a,
		}
	}
}

func (r *ErrorResponse) populateErrorRaw() {
	raw := r.ErrorRaw
	var err error
	var s string
	if err = json.Unmarshal(raw, &s); err == nil {
		r.Message = s
		return
	}
	var m map[string][]string
	if err = json.Unmarshal(raw, &m); err == nil {
		r.ValidationError = m
	}
}

func (m errorMap) String() string {
	var buffer bytes.Buffer
	for k, errors := range m {
		s := fmt.Sprintf("[%s: %s]", k, strings.Join(errors, ", "))
		buffer.WriteString(s)
	}
	return buffer.String()
}

func (r *ErrorResponse) Error() string {
	m := fmt.Sprintf("%v %v: %d",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode)
	if r.Message != "" {
		m += " " + r.Message
	}
	if len(r.ValidationError) > 0 {
		m += " " + r.ValidationError.String()
	}
	return m
}
