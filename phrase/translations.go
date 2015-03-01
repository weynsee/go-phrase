package phrase

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TranslationsService struct {
	client *Client
}

type Translation struct {
	ID                 int      `json:"id" url:"-"`
	Content            string   `json:"content" url:"content"`
	PluralSuffix       string   `json:"plural_suffix" url:"plural_suffix"`
	Placeholders       []string `json:"placeholders" url:"-"`
	Unverified         bool     `json:"unverified" url:"-"`
	ExcludedFromExport bool     `json:"excluded_from_export" url:"excluded_from_export,int,omitempty"`
	Key                Key      `json:"translation_key" url:"-"`
}

type DownloadRequest struct {
	Locale                   string    `url:"locale"`
	Format                   string    `url:"format"`
	UpdatedSince             time.Time `url:"updated_since,omitempty"`
	Tag                      string    `url:"tag,omitempty"`
	IncludeEmptyTranslations bool      `url:"include_empty_translations,int,omitempty"`
	KeepNotranslateTags      bool      `url:"keep_notranslate_tags,int,omitempty"`
	ConvertEmoji             bool      `url:"convert_emoji,int,omitempty"`
}

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

const timeFormat = "20060102150405"

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

func (r RateLimit) String() string {
	return fmt.Sprintf("RateLimit Limit: %d Remaining: %d Reset: %v",
		r.Limit, r.Remaining, r.Reset)
}

func (t Translation) String() string {
	return fmt.Sprintf("Translation ID: %d Content: %s Key: %v",
		t.ID, t.Content, t.Key)
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
