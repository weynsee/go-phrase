package phrase

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://phraseapp.com/api/v1/"
	userAgent      = "go-phrase/" + libraryVersion
)

type Client struct {
	client    *http.Client
	baseURL   *url.URL
	UserAgent string

	AuthToken        string
	ProjectAuthToken string

	Sessions     *SessionsService
	Projects     *ProjectsService
	Locales      *LocalesService
	Keys         *KeysService
	Blacklist    *BlacklistService
	Tags         *TagsService
	Translations *TranslationsService
	FileImports  *FileImportsService
	Orders       *OrdersService
}

func New(authtoken string) *Client {
	return NewClient(authtoken, "", nil)
}

func NewClient(authToken, projectToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{AuthToken: authToken, ProjectAuthToken: projectToken,
		client: httpClient, baseURL: baseURL, UserAgent: userAgent}
	c.Sessions = &SessionsService{c}
	c.Projects = &ProjectsService{c}
	c.Locales = &LocalesService{c}
	c.Keys = &KeysService{c}
	c.Blacklist = &BlacklistService{c}
	c.Tags = &TagsService{c}
	c.Translations = &TranslationsService{c}
	c.FileImports = &FileImportsService{c}
	c.Orders = &OrdersService{c}
	return c
}

func (c *Client) NewRequest(method, urlStr string, params url.Values) (*http.Request, error) {
	if params == nil {
		params = url.Values{}
	}
	c.setDefaultParams(&params)
	u, err := c.resolveURL(urlStr)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	method = strings.ToUpper(method)
	if method == "GET" {
		u.RawQuery = params.Encode()
	} else {
		body = strings.NewReader(params.Encode())
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if method != "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}
	return req, nil
}

func (c *Client) setDefaultParams(params *url.Values) {
	params.Set("auth_token", c.AuthToken)
	if c.ProjectAuthToken != "" {
		params.Set("project_auth_token", c.ProjectAuthToken)
	}
}

func (c *Client) resolveURL(urlStr string) (*url.URL, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return c.baseURL.ResolveReference(rel), nil
}

func (c *Client) NewUploadRequest(urlStr string, params url.Values, paramName, filename string, reader io.Reader) (*http.Request, error) {
	if params == nil {
		params = url.Values{}
	}
	c.setDefaultParams(&params)

	u, err := c.resolveURL(urlStr)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, reader)

	for key, _ := range params {
		_ = writer.WriteField(key, params.Get(key))
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", u.String(), body)

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, err
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return resp, err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	return ResponseError(r)
}

type successResponse struct {
	Success bool `json:"success"`
}
