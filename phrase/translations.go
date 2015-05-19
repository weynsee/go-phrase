package phrase

import (
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// TranslationsService provides access to the translation related functions
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translations/
type TranslationsService struct {
	client *Client
}

// Translation represents a translation stored in PhraseApp.
type Translation struct {
	ID      int    `json:"id" url:"-"`
	Content string `json:"content" url:"content"`

	// A plural suffix (only required when the key is pluralized).
	PluralSuffix string   `json:"plural_suffix" url:"plural_suffix"`
	Placeholders []string `json:"placeholders" url:"-"`
	Unverified   bool     `json:"unverified" url:"-"`

	// Whether the translation should be excluded from file downloads.
	ExcludedFromExport bool `json:"excluded_from_export" url:"excluded_from_export,int,omitempty"`
	Key                Key  `json:"translation_key" url:"-"`
}

// DownloadRequest represents the parameters to the download API call.
type DownloadRequest struct {
	// Name of the locale that should be downloaded. This field is mandatory.
	Locale string `url:"locale"`

	// Name of the format that should be downloaded. This field is mandatory.
	// http://docs.phraseapp.com/guides/formats/
	Format string `url:"format"`

	// Encoding of the downloaded file
	Encoding string `url:"encoding,omitempty"`

	// Date since the last update. Only return translations that were modified after the given date.
	UpdatedSince time.Time `url:"updated_since,omitempty"`
	Tag          string    `url:"tag,omitempty"`

	// Include empty translations in the result as well.
	IncludeEmptyTranslations bool `url:"include_empty_translations,int,omitempty"`

	// Do not remove NOTRANSLATE tags in the downloaded translations.
	KeepNotranslateTags bool `url:"keep_notranslate_tags,int,omitempty"`

	// Enable Emoji conversion.
	ConvertEmoji bool `url:"convert_emoji,int,omitempty"`

	// Skip unverified translations from appearing in the downloaded file.
	SkipUnverifiedTranslations bool `url:"skip_unverified_translations,int,omitempty"`
}

// RateLimit represents the rate limits returned in the response.
type RateLimit struct {
	// Number of max requests allowed in the current time period.
	Limit int

	// Number of remaining requests in the current time period.
	Remaining int

	// Timestamp of end of current time period as UNIX timestamp.
	Reset time.Time
}

const timeFormat = "20060102150405"

// List all translations for a locale.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translations/#index
func (s *TranslationsService) Get(l string, t *time.Time) ([]Translation, error) {
	params := url.Values{}
	params.Set("locale_name", l)
	if t != nil {
		params.Set("updated_since", t.Format(timeFormat))
	}
	req, err := s.client.NewRequest("GET", "translations", params)
	if err != nil {
		return nil, err
	}

	var translations []Translation
	_, err = s.client.Do(req, &translations)
	if err != nil {
		return nil, err
	}

	return translations, err
}

// List all translations for all locale.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translations/#index
func (s *TranslationsService) ListAll() (map[string][]Translation, error) {
	req, err := s.client.NewRequest("GET", "translations", nil)
	if err != nil {
		return nil, err
	}

	var translations map[string][]Translation
	_, err = s.client.Do(req, &translations)
	if err != nil {
		return nil, err
	}

	return translations, err
}

// Fetch translations for a locale and a list of keys.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translations/#fetch_list
func (s *TranslationsService) GetByKeys(l string, keys []string) ([]Translation, error) {
	params := url.Values{}
	params.Set("locale", l)
	for _, key := range keys {
		params.Add("keys[]", key)
	}
	req, err := s.client.NewRequest("POST", "translations/fetch_list", params)
	if err != nil {
		return nil, err
	}

	var translations []Translation
	_, err = s.client.Do(req, &translations)
	if err != nil {
		return nil, err
	}

	return translations, err
}

// Download translations as localization file. This call is rate limited.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translations/#download
func (s *TranslationsService) Download(d *DownloadRequest, w io.Writer) (*RateLimit, error) {
	params, err := query.Values(d)
	if err != nil {
		return nil, err
	}
	if !d.UpdatedSince.IsZero() {
		params.Set("updated_since", d.UpdatedSince.Format(timeFormat))
	}
	req, err := s.client.NewRequest("GET", "translations/download", params)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req, w)
	rate := getRateLimit(&resp.Header)
	if err != nil {
		return rate, err
	}
	return rate, nil
}

// Save a translation for a given locale.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translations/#store
func (s *TranslationsService) Update(locale, key string, t *Translation, skipVerification, disallowUpdate bool) (*Translation, error) {
	params, err := query.Values(t)
	params.Set("locale", locale)
	params.Set("key", key)
	if skipVerification {
		params.Set("skip_verification", "1")
	}
	if disallowUpdate {
		params.Set("allow_update", "0")
	}
	req, err := s.client.NewRequest("POST", "translations/store", params)
	if err != nil {
		return nil, err
	}

	translation := new(Translation)
	_, err = s.client.Do(req, translation)
	if err != nil {
		return nil, err
	}

	return translation, err
}

func getRateLimit(h *http.Header) *RateLimit {
	rate := new(RateLimit)
	if limit := h.Get("X-Rate-Limit-Limit"); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := h.Get("X-Rate-Limit-Remaining"); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := h.Get("X-Rate-Limit-Reset"); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.Reset = time.Unix(v, 0)
		}
	}
	return rate
}
