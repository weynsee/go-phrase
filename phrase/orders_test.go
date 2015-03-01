package phrase

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrdersService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_orders", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"amount_in_cents":1152,"lsp":"gengo","code":"30AB4884","currency":"usd","message":"","state":"confirmed","translation_type":"pro","progress_percent":0,"source_locale_name":"en","target_locale_names":["de"],"source_locale_code":"en-GB","target_locale_codes":["de-DE"],"tag_name":null,"styleguide_url":null,"styleguide":null},{"amount_in_cents":960,"lsp":"textmaster","code":"046037B2","currency":"usd","message":"","state":"open","translation_type":"pro","progress_percent":0,"source_locale_name":"en","source_locale_code":"en-GB","target_locale_names":["de"],"target_locale_codes":["de-DE"],"tag_name":"foo-tag","styleguide_url":null,"styleguide":null}]`)
	})

	orders, err := client.Orders.ListAll()
	if err != nil {
		t.Errorf("Orders.ListAll returned error: %v", err)
	}

	want := []Order{
		{
			AmountInCents: 1152, LSP: "gengo",
			Code: "30AB4884", Currency: "usd", State: "confirmed",
			TranslationType: "pro", ProgressPercent: 0,
			SourceLocaleName: "en", TargetLocaleNames: []string{"de"},
			SourceLocaleCode: "en-GB", TargetLocaleCodes: []string{"de-DE"},
		},
		{
			AmountInCents: 960, LSP: "textmaster",
			Code: "046037B2", Currency: "usd", State: "open",
			TranslationType: "pro", ProgressPercent: 0,
			SourceLocaleName: "en", TargetLocaleNames: []string{"de"},
			SourceLocaleCode: "en-GB", TargetLocaleCodes: []string{"de-DE"},
			Tag: "foo-tag",
		},
	}
	if !reflect.DeepEqual(orders, want) {
		t.Errorf("Orders.ListAll returned %+v, want %+v", orders, want)
	}
}

func TestOrdersService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_orders/CODE", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"amount_in_cents":1152,"lsp":"gengo","code":"30AB4884","currency":"usd","message":"","state":"confirmed","translation_type":"pro","progress_percent":0,"source_locale_name":"en","target_locale_names":["de"],"source_locale_code":"en-GB","target_locale_codes":["de-DE"],"tag_name":null,"styleguide_url":null,"styleguide":null}`)
	})

	order, err := client.Orders.Get("CODE")
	if err != nil {
		t.Errorf("Orders.Get returned error: %v", err)
	}

	want := &Order{
		AmountInCents: 1152, LSP: "gengo",
		Code: "30AB4884", Currency: "usd", State: "confirmed",
		TranslationType: "pro", ProgressPercent: 0,
		SourceLocaleName: "en", TargetLocaleNames: []string{"de"},
		SourceLocaleCode: "en-GB", TargetLocaleCodes: []string{"de-DE"},
	}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("Orders.Get returned %+v, want %+v", order, want)
	}
}

func TestOrdersService_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_orders/this", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Orders.Destroy("this")
	if err != nil {
		t.Errorf("Orders.Destroy returned error: %v", err)
	}
}

func TestOrdersService_Confirm(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_orders/30AB4884/confirm", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"amount_in_cents":1152,"lsp":"gengo","code":"30AB4884","currency":"usd","message":"Hello World","state":"confirmed","translation_type":"pro","progress_percent":0,"source_locale_name":"en","target_locale_names":["de"],"tag_name":null,"styleguide_url":null,"unverify_translations_upon_delivery":true,"quality":true,"priority":true,"expertise":true,"styleguide":null}`)
	})

	order, err := client.Orders.Confirm("30AB4884")
	if err != nil {
		t.Errorf("Orders.Confirm returned error: %v", err)
	}

	want := &Order{
		AmountInCents:                    1152,
		LSP:                              "gengo",
		Code:                             "30AB4884",
		Currency:                         "usd",
		Message:                          "Hello World",
		State:                            "confirmed",
		TranslationType:                  "pro",
		ProgressPercent:                  0,
		SourceLocaleName:                 "en",
		TargetLocaleNames:                []string{"de"},
		UnverifyTranslationsUponDelivery: true,
		Quality:   true,
		Priority:  true,
		Expertise: true,
	}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("Orders.Confirm returned %+v, want %+v", order, want)
	}
}

func TestOrdersService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_orders", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"source_locale_name":   "en",
			"target_locale_name[]": "fr",
			"translation_type":     "pro",
		})
		fmt.Fprint(w, `{"amount_in_cents":1152,"lsp":"textmaster","code":"30AB4884","currency":"usd","message":"Hello World","state":"open","translation_type":"pro","progress_percent":0,"source_locale_name":"en","target_locale_names":["de"],"tag_name":null,"styleguide_url":null,"unverify_translations_upon_delivery":true,"quality":true,"priority":true,"expertise":true,"styleguide":null}`)
	})

	order, err := client.Orders.Create(
		&Order{
			LSP:               "gengo",
			SourceLocaleName:  "en",
			TargetLocaleNames: []string{"fr"},
			TranslationType:   "pro",
		},
	)
	if err != nil {
		t.Errorf("Orders.Create returned error: %v", err)
	}

	want := &Order{
		AmountInCents:                    1152,
		LSP:                              "textmaster",
		Code:                             "30AB4884",
		Currency:                         "usd",
		Message:                          "Hello World",
		State:                            "open",
		TranslationType:                  "pro",
		ProgressPercent:                  0,
		SourceLocaleName:                 "en",
		TargetLocaleNames:                []string{"de"},
		UnverifyTranslationsUponDelivery: true,
		Quality:   true,
		Priority:  true,
		Expertise: true,
	}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("Orders.Create returned %+v, want %+v", order, want)
	}
}
