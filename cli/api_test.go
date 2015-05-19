package cli

import (
	"github.com/weynsee/go-phrase/phrase"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

var (
	mux    *http.ServeMux
	client *phrase.Client
	server *httptest.Server

	testFolder = "test"
)

func setupAPI() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = phrase.New("faketoken")
	client.BaseURL, _ = url.Parse(server.URL)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
}

func tearDown() {
	shutdownAPI()
	os.RemoveAll(testFolder)
}

func shutdownAPI() {
	server.Close()
}
