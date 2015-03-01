package phrase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestResponseError_messageText(t *testing.T) {
	message := "Bad Request"
	resp := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(message)),
	}
	if e := ResponseError(resp); e.Message != message {
		t.Errorf("ErrorResponse Message as String = %v, want %v", e.Message, message)
	}
}

func TestResponseError_string(t *testing.T) {
	url, _ := url.Parse("http://www.example.com")
	req := &http.Request{
		Method: "PUT",
		URL:    url,
	}
	resp := &http.Response{
		Request:    req,
		StatusCode: http.StatusBadRequest,
	}
	message := "validation error"
	validation := map[string][]string{
		"field1": []string{"is wrong", "is incorrect"},
	}
	ep := ErrorResponse{
		Response:        resp,
		Message:         message,
		ValidationError: validation,
	}
	want := "PUT http://www.example.com: 400 validation error [field1: is wrong, is incorrect]"
	if got := ep.Error(); got != want {
		t.Errorf("ErrorResponse Error() %v, want %v", got, want)
	}
}

func TestResponseError_messageJSONString(t *testing.T) {
	validateErrorResponseBodyString(t, "messages")
}

func TestResponseError_errorJSONString(t *testing.T) {
	validateErrorResponseBodyString(t, "error")
}

func validateErrorResponseBodyString(t *testing.T, key string) {
	message := "not found"
	json := fmt.Sprintf(`{"%v":"%s"}`, key, message)
	if e := ResponseError(jsonErrorResponse(json, http.StatusNotFound)); e.Message != message {
		t.Errorf("ErrorResponse Message as JSON String = %v, want %v", e.Message, message)
	}
}

func TestResponseError_messageJSONArray(t *testing.T) {
	messages := []string{"first error", "second error"}
	json := fmt.Sprintf(`{"messages":["%s", "%s"]}`, messages[0], messages[1])
	resp := jsonErrorResponse(json, 422)
	e := ResponseError(resp)
	validateValidationError(t, e.ValidationError, "error", messages)
}

func validateValidationError(t *testing.T, em errorMap, key string, errors []string) {
	a, exists := em[key]
	if !exists {
		t.Fatalf("ErrorResponse ValidationError %v key not found", key)
	}
	if got, want := len(a), len(errors); got != want {
		t.Fatalf("ErrorResponse ValidationError length = %d, want %d", got, want)
	}
	for i, message := range errors {
		if a[i] != message {
			t.Errorf("ErrorResponse ValidationError message = %v, want %v", a[i], message)
		}
	}
}

func TestResponseError_messageJSONMap(t *testing.T) {
	validateErrorResponseBody(t, "messages")
}

func TestResponseError_errorJSONMap(t *testing.T) {
	validateErrorResponseBody(t, "error")
}

func validateErrorResponseBody(t *testing.T, key string) {
	error1, error2 := []string{"first error", "second error"}, []string{"third error"}
	messages := map[string][]string{
		"field1": error1,
		"field2": error2,
	}
	body, _ := json.Marshal(map[string]errorMap{key: messages})
	resp := jsonErrorResponse(string(body), 422)
	e := ResponseError(resp)
	validateValidationError(t, e.ValidationError, "field1", error1)
	validateValidationError(t, e.ValidationError, "field2", error2)
}

func jsonErrorResponse(body string, status int) *http.Response {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return &http.Response{
		Request:    &http.Request{},
		StatusCode: status,
		Header:     header,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
}
