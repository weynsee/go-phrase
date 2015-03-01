package phrase

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = New("faketoken")
	client.baseURL, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestNew(t *testing.T) {
	c := New("token1")
	if got, want := c.AuthToken, "token1"; got != want {
		t.Errorf("New AuthToken is %v, want %v", got, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("token1", "token2", nil)

	if got, want := c.baseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient baseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}
	if got, want := c.AuthToken, "token1"; got != want {
		t.Errorf("NewClient AuthToken is %v, want %v", got, want)
	}
	if got, want := c.ProjectAuthToken, "token2"; got != want {
		t.Errorf("NewClient ProjectAuthToken is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient("token1", "token2", nil)
	inURL, outURL := "test", defaultBaseURL+"test?auth_token=token1&project_auth_token=token2"
	req, _ := c.NewRequest("GET", inURL, nil)
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestNewRequest_GET(t *testing.T) {
	c := New("token")
	inURL, outURL := "test", defaultBaseURL+"test?auth_token=token&key%5B%5D=key1"
	params := url.Values{}
	params.Add("key[]", "key1")
	req, _ := c.NewRequest("GET", inURL, params)
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}
	if got, want := req.Method, "GET"; got != want {
		t.Errorf("NewRequest() Method is %v, want %v", got, want)
	}
}

func TestNewRequest_POST(t *testing.T) {
	c := New("token")
	inURL, outURL := "test", defaultBaseURL+"test"
	params := url.Values{}
	params.Add("post_param", "value")
	req, _ := c.NewRequest("POST", inURL, params)
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), "auth_token=token&post_param=value"; got != want {
		t.Errorf("NewRequest() Body is %v, want %v", got, want)
	}
	if got, want := req.Method, "POST"; got != want {
		t.Errorf("NewRequest() Method is %v, want %v", got, want)
	}
	if got, want := req.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; got != want {
		t.Errorf("NewRequest() Content-Type is %v, want %v", got, want)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := New("")
	_, err := c.NewRequest("GET", ":", nil)
	testParseURLError(t, err)
}

func testParseURLError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewUploadRequest(t *testing.T) {
	c := NewClient("token1", "token2", nil)
	params := url.Values{}
	params.Set("param", "value")
	contents := "this is a test file upload"
	file := strings.NewReader(contents)
	req, _ := c.NewUploadRequest("test", params, "file", "filename.txt", file)
	values := map[string]string{
		"param":              "value",
		"auth_token":         "token1",
		"project_auth_token": "token2",
	}
	body, _ := ioutil.ReadAll(req.Body)
	resp := string(body)
	for k, v := range values {
		val := fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\r\n\r\n%s", k, v)
		if i := strings.Index(resp, val); i == -1 {
			t.Errorf("NewUploadRequest() did not encode the value %v in the body", v)
		}
	}
	f := fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"; filename=\"%s\"\r\n", "file", "filename.txt")
	if i := strings.Index(resp, f); i == -1 {
		t.Error("NewUploadRequest() did not encode the filename and file param in the body")
	}
	if i := strings.Index(resp, contents); i == -1 {
		t.Error("NewUploadRequest() did not encode the file in the body")
	}
	if got, want := req.Method, "POST"; got != want {
		t.Errorf("NewUploadRequest() Method is %v, want %v", got, want)
	}
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewUploadRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestNewUploadRequest_badURL(t *testing.T) {
	c := New("")
	_, err := c.NewUploadRequest(":", nil, "", "", nil)
	testParseURLError(t, err)
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error.")
	}
}

func testFormValues(t *testing.T, r *http.Request, want map[string]string) {
	body, _ := ioutil.ReadAll(r.Body)
	values, _ := url.ParseQuery(string(body))
	testParams(t, values, want)
}

func testQueryString(t *testing.T, r *http.Request, want map[string]string) {
	values, _ := url.ParseQuery(r.URL.RawQuery)
	testParams(t, values, want)
}

func testParams(t *testing.T, values url.Values, want map[string]string) {
	for k, v := range want {
		if values.Get(k) != v {
			t.Errorf("Request body parameter %s should have value %s", k, v)
		}
	}
}
