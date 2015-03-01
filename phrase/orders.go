package phrase

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
)

type OrdersService struct {
	client *Client
}

type Order struct {
	AmountInCents                    int      `json:"amount_in_cents" url:"-"`
	LSP                              string   `json:"lsp" url:"lsp"`
	Code                             string   `json:"code" url:"-"`
	Currency                         string   `json:"currency" url:"-"`
	Message                          string   `json:"message" url:"message,omitempty"`
	State                            string   `json:"state" url:"-"`
	TranslationType                  string   `json:"translation_type" url:"translation_type"`
	ProgressPercent                  int      `json:"progress_percent" url:"-"`
	SourceLocaleName                 string   `json:"source_locale_name" url:"source_locale_name"`
	SourceLocaleCode                 string   `json:"source_locale_code" url:"-"`
	TargetLocaleNames                []string `json:"target_locale_names" url:"target_locale_name[]"`
	TargetLocaleCodes                []string `json:"target_locale_codes" url:"-"`
	Tag                              string   `json:"tag_name" url:"tag_name,omitempty"`
	StyleguideUrl                    string   `json:"styleguide_url" url:"-"`
	Styleguide                       string   `json:"styleguide" url:"-"`
	StyleguideCode                   string   `json:"styleguide_code" url:"styleguide_code,omitempty"`
	UnverifyTranslationsUponDelivery bool     `json:"unverify_translations_upon_delivery" url:"unverify_translations_upon_delivery,int,omitempty"`
	Quality                          bool     `json:"quality" url:"quality,int,omitempty"`
	Priority                         bool     `json:"priority" url:"priority,int,omitempty"`
	Expertise                        bool     `json:"expertise" url:"expertise,int,omitempty"`
	Category                         int      `json:"category" url:"category,omitempty"`
}

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

func (s *OrdersService) Create(o *Order) (*Order, error) {
	params, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	return s.submitOrder("POST", "translation_orders", params)
}

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

func (o Order) String() string {
	return fmt.Sprintf("Order Code: %s LSP: %s State: %s Amount in Cents: %d",
		o.Code, o.LSP, o.State, o.AmountInCents)
}
