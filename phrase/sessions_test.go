package phrase

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSessionsService_Create(t *testing.T) {
	setup()
	defer teardown()

	email, password, token := "email@mail.com", "password", "sometoken"

	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"email":    email,
			"password": password,
		})
		fmt.Fprint(w, fmt.Sprintf(`{"success":true,"auth_token":"%s"}`, token))
	})

	got, err := client.Sessions.Create(email, password)
	if err != nil {
		t.Errorf("Sessions.Create returned error: %v", err)
	}

	if token != got {
		t.Errorf("Sessions.Create returned %+v, want %+v", got, token)
	}
}

func TestSessionsService_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"success":true}`)
	})

	err := client.Sessions.Destroy()
	if err != nil {
		t.Errorf("Sessions.Destroy returned error: %v", err)
	}
}

func TestSessionsService_CheckLogin(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/auth/check_login", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"logged_in":true,"user":{"id":1,"email":"demo@phraseapp.com"}}`)
	})

	user, err := client.Sessions.CheckLogin()
	if err != nil {
		t.Errorf("Sessions.CheckLogin returned error: %v", err)
	}

	want := &User{ID: 1, Email: "demo@phraseapp.com"}
	if !reflect.DeepEqual(user, want) {
		t.Errorf("Sessions.CheckLogin returned %+v, want %+v", user, want)
	}
}
