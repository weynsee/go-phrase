package phrase

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// LocalesService provides access to the locales related functions
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/locales/
type LocalesService struct {
	client *Client
}

// Locale represents a locale.
type Locale struct {
	ID             int                          `json:"id"`
	Name           string                       `json:"name"`
	Code           string                       `json:"code"`
	CountryCode    string                       `json:"country_code"`
	Direction      string                       `json:"writing_direction"`
	Default        bool                         `json:"is_default"`
	Pluralizations map[string]map[string]string `json:"pluralizations"`
}

// Returns list of existing locales in current project.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/locales/#index
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

// Download translations for locale in a specific format.
// See http://docs.phraseapp.com/guides/formats/ for the list of all supported formats.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/locales/#show
func (s *LocalesService) Download(locale, format string, w io.Writer) error {
	path := fmt.Sprintf("locales/%s.%s", locale, format)
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

// Create a new locale in the current project.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/locales/#create
func (s *LocalesService) Create(name string) (*Locale, error) {
	params := url.Values{}
	params.Set("locale[name]", name)
	req, err := s.client.NewRequest("POST", "locales", params)
	if err != nil {
		return nil, err
	}

	return s.requestLocale(req)
}

// Promotes locale to be the default locale for the current project.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/locales/#make_default
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
