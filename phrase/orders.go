package phrase

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
)

// OrdersService provides access to the translation order related functions
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translation_orders/
type OrdersService struct {
	client *Client
}

// Order represents a translation order.
type Order struct {
	AmountInCents int `json:"amount_in_cents" url:"-"`

	// Name of the LSP you want to use for this order (can be "gengo" or "textmaster").
	LSP      string `json:"lsp" url:"lsp"`
	Code     string `json:"code" url:"-"`
	Currency string `json:"currency" url:"-"`

	// Give an optional message to be delivered to the translators.
	Message string `json:"message" url:"message,omitempty"`
	State   string `json:"state" url:"-"`

	// Quality level of the translations. Can be "standard" or "pro" for Gengo orders and "regular" or "premium" for TextMaster orders.
	TranslationType string `json:"translation_type" url:"translation_type"`
	ProgressPercent int    `json:"progress_percent" url:"-"`

	// Name of the source locale. This locale must already exist in your project. This usually is your default locale.
	SourceLocaleName string `json:"source_locale_name" url:"source_locale_name"`
	SourceLocaleCode string `json:"source_locale_code" url:"-"`

	// Name of the locale to translate your content into. This locale must already exist in your project.
	TargetLocaleNames []string `json:"target_locale_names" url:"target_locale_name[]"`
	TargetLocaleCodes []string `json:"target_locale_codes" url:"-"`

	// Name of a tag to limit the translations to. If none given, it is assumed that all existing keys should be translated.
	Tag           string `json:"tag_name" url:"tag_name,omitempty"`
	StyleguideURL string `json:"styleguide_url" url:"-"`
	Styleguide    string `json:"styleguide" url:"-"`

	// Identification code of the style guide to attach with this order.
	StyleguideCode string `json:"styleguide_code" url:"styleguide_code,omitempty"`

	// Enable unverifying translations upon delivery.
	UnverifyTranslationsUponDelivery bool `json:"unverify_translations_upon_delivery" url:"unverify_translations_upon_delivery,int,omitempty"`

	// Use this option to order translations for keys with unverified content in the selected target locales.
	IncludeUnverifiedTranslations bool `json:"include_unverified_translations" url:"include_unverified_translations,int,omitempty"`

	// Use this option to order translations for keys with untranslated content in the selected target locales.
	IncludeUntranslatedKeys bool `json:"include_untranslated_keys" url:"include_untranslated_keys,int,omitempty"`

	// Translations will be proofread by TextMaster to ensure consistency in vocabulary and style. Will cause additional costs! Only available for TextMaster!
	Quality bool `json:"quality" url:"quality,int,omitempty"`

	// Your project will be assigned a higher priority status, which decreases turnaround time by 30%. Will cause additional costs! Only available for TextMaster!
	Priority bool `json:"priority" url:"priority,int,omitempty"`

	// TextMaster provides you with an expert in the selected category. Will cause additional costs! Only available for TextMaster!
	Expertise bool `json:"expertise" url:"expertise,int,omitempty"`

	// Category ID to use (only required for orders through TextMaster).
	// http://docs.phraseapp.com/api/v1/translation_orders/#categories
	Category int `json:"category" url:"category,omitempty"`
}

// ListAll gets a list of all orders in the current project.
// This is a signed request.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translation_orders/#index
func (s *OrdersService) ListAll() ([]Order, error) {
	req, err := s.client.NewRequest("GET", "translation_orders", nil)
	if err != nil {
		return nil, err
	}

	orders := new([]Order)
	_, err = s.client.Do(req, orders)
	if err != nil {
		return nil, err
	}

	return *orders, err
}

// Get details of the order identified by the order code.
// This is a signed request.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translation_orders/#show
func (s *OrdersService) Get(code string) (*Order, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("translation_orders/%s", code), nil)
	if err != nil {
		return nil, err
	}

	order := new(Order)
	_, err = s.client.Do(req, order)
	if err != nil {
		return nil, err
	}

	return order, err
}

// Create a new order to confirm.
// This is a signed request.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translation_orders/#create
func (s *OrdersService) Create(o *Order) (*Order, error) {
	params, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	return s.submitOrder("POST", "translation_orders", params)
}

// Destroy deletes an order (must not yet be confirmed).
// This is a signed request.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translation_orders/#destroy
func (s *OrdersService) Destroy(code string) error {
	u := fmt.Sprintf("translation_orders/%s", code)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return err
}

// Confirm confirms an open order identified by the order code and starts the translation process. Valid billing information is required for this action. After confirming, the displayed amount will be charged.
// This is a signed request.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/translation_orders/#confirm
func (s *OrdersService) Confirm(code string) (*Order, error) {
	return s.submitOrder("PUT", fmt.Sprintf("translation_orders/%s/confirm", code), nil)
}

func (s *OrdersService) submitOrder(method, url string, params url.Values) (*Order, error) {
	req, err := s.client.NewRequest(method, url, params)
	if err != nil {
		return nil, err
	}

	order := new(Order)
	_, err = s.client.Do(req, order)
	if err != nil {
		return nil, err
	}

	return order, err
}
