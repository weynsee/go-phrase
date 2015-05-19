package phrase

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBlacklistService_Keys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/blacklisted_keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"name":"key1"},{"name":"key2"}]`)
	})

	keys, err := client.Blacklist.Keys()
	if err != nil {
		t.Errorf("Blacklist.Keys returned error: %v", err)
	}

	want := []string{"key1", "key2"}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("Blacklist.Keys returned %+v, want %+v", keys, want)
	}
}

func TestBlacklistService_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Blacklist.Keys()
		return err
	})
}
