package phrase

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type LocalesService struct {
	client *Client
}

type Locale struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	CountryCode string `json:"country_code"`
	Direction   string `json:"writing_direction"`
	Default     bool   `json:"is_default"`
}

func (s *LocalesService) ListAll() ([]Locale, error) {
	req, err := s.client.NewRequest("GET", "locales", nil)
	if err != nil {
		return nil, err
	}

	locales := new([]Locale)
	_, err = s.client.Do(req, locales)
	if err != nil {
		return nil, err
	}

	return *locales, err
}

func (l Locale) String() string {
	return fmt.Sprintf("Locale ID: %d Name: %s",
		l.ID, l.Name)
}

func (s *LocalesService) Download(locale, f string, w io.Writer) error {
	path := fmt.Sprintf("locales/%s.%s", locale, f)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, w)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalesService) Create(name string) (*Locale, error) {
	params := url.Values{}
	params.Set("locale[name]", name)
	req, err := s.client.NewRequest("POST", "locales", params)
	if err != nil {
		return nil, err
	}

	return s.requestLocale(req)
}

func (s *LocalesService) MakeDefault(name string) (*Locale, error) {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("locales/%s/make_default", name), nil)
	if err != nil {
		return nil, err
	}

	return s.requestLocale(req)
}

func (s *LocalesService) requestLocale(req *http.Request) (*Locale, error) {
	locale := new(Locale)
	_, err := s.client.Do(req, locale)
	if err != nil {
		return nil, err
	}

	return locale, err
}
