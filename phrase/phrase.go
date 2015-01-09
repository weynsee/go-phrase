package phrase

import (
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://phraseapp.com/api/v1"
	userAgent      = "go-phrase/" + libraryVersion
)

type Client struct {
	client    *http.Client
	baseURL   *url.URL
	UserAgent string

	AuthToken string

	Projects *ProjectsService
}

func NewClient(authToken, projectToken string, httpClient *http.Client) *Client {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, baseURL: baseURL, UserAgent: userAgent}
	c.Projects = &ProjectsService{c}
	return c
}
